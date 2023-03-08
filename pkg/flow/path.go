package flow

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	ErrStopGoroutine    = errors.New("stop goroutine")
	ErrEndpointNotFound = errors.New("endpoint not found")
)

type (
	ContextTagsValue = map[string]struct{}
	ContextType      string
)

const CtxTags ContextType = "tags"

func VisitAndFetch(ctx context.Context, reg *NodesReg) error {
	starts := reg.GetStarts()

	if starts == nil {
		return fmt.Errorf("%s %w", reg.startName, ErrEndpointNotFound)
	}

	// gather tags for same start
	// TODO: this is not the best way to do this
	tags := make(map[string]struct{})
	for _, start := range starts {
		for tag := range start.Tags {
			tags[tag] = struct{}{}
		}
	}

	ctx = context.WithValue(ctx, CtxTags, tags)

	for _, start := range starts {
		// add tags to context
		if err := validateFetch(ctx, "", []Connection{start.Connection}, reg); err != nil {
			return err
		}
	}

	return nil
}

func validateFetch(ctx context.Context, current string, outputs []Connection, reg *NodesReg) error {
	// log.Debug().Msgf("current %s", current)
	for _, output := range outputs {
		// log.Debug().Msgf("validating [%v]", output)
		// before to run check context
		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:
		}

		node, ok := reg.Get(output.Node)
		if !ok {
			return fmt.Errorf("node not found %s", output.Node)
		}

		// is disabled node
		if node.IsDisabled() {
			continue
		}

		// activate this input
		ctxTags, _ := ctx.Value(CtxTags).(ContextTagsValue)
		node.ActiveInput(current, ctxTags)

		if node.IsChecked() {
			continue
		}

		// fetch and validation

		if err := node.Validate(ctx); err != nil {
			return fmt.Errorf("ID %s, %s validate failed: %v", output, node.GetType(), err)
		}

		if err := node.Fetch(ctx, reg.appStore.DB); err != nil {
			return fmt.Errorf("ID %s, %s fetch failed: %v", output, node.GetType(), err)
		}

		// respond channel activate
		if node.IsRespond() && !reg.respondChanActive {
			reg.respondChan = make(chan Respond, 1)
			reg.respondChanActive = true
		}

		// stamp check
		node.Check()

		for i := 0; i < node.NextCount(); i++ {
			if err := validateFetch(ctx, output.Node, node.Next(i), reg); err != nil {
				return err
			}
		}
	}

	return nil
}

func GoAndRun(ctx context.Context, wg *sync.WaitGroup, reg *NodesReg, firstValue []byte) {
	defer wg.Done()

	starts := reg.GetStarts()

	// stuct count check

	reg.stuckChan = make(chan bool, 1)

	stuckCheckCtx, stuckCheckCtxCancel := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			// everytime a job is finished
			case cancel := <-reg.stuckChan:
				if cancel {
					// all job is finished

					reg.CancelStucks()
				}
			case <-stuckCheckCtx.Done():
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	// change waitgroup to check all job is finished
	for _, start := range starts {
		branch(ctx, []Connection{start.Connection}, reg, &nodeRetOutput{firstValue})
	}

	reg.wgx.Wait()

	// cancel stuck check
	stuckCheckCtxCancel()

	if reg.respondChan != nil {
		if reg.respondChanActive {
			// send error to channel
			buff := bytes.Buffer{}
			for _, err := range reg.errors {
				buff.WriteString("[")
				buff.WriteString(err.Error())
				buff.WriteString("]")
			}

			if buff.Len() != 0 {
				reg.respondChan <- Respond{
					Data:    buff.Bytes(),
					IsError: true,
				}
			} else {
				reg.respondChan <- Respond{
					Status:  http.StatusAccepted,
					Data:    []byte("Accepted"),
					IsError: false,
				}
			}
		}

		close(reg.respondChan)
	}

	log.Ctx(ctx).Info().Msgf("completed control flow")
}

// branch values already designed to same length of nexts.
func branch(ctx context.Context, nexts []Connection, reg *NodesReg, value NodeRet) {
	for _, next := range nexts {
		// going goroutine to prevent too much recursive call
		reg.wgx.Add(1)
		reg.UpdateStuck(CountTotalIncrease, false)

		go branchRun(ctx, next, reg, value)
	}
}

func branchRun(ctx context.Context, start Connection, reg *NodesReg, value NodeRet) {
	defer func() {
		// check panic
		if r := recover(); r != nil {
			log.Ctx(ctx).Error().Msgf("panic: %v\n%v", r, string(debug.Stack()))
			reg.AddError(fmt.Errorf("panic: %s cannot run: %v\n%v", start.Node, r, string(debug.Stack())))
		}

		reg.UpdateStuck(CountTotalDecrease, true)
		reg.wgx.Done()
	}()

	// defer reg.wgx.Done()
	// defer reg.UpdateStuck(CountTotalDecrease, true)

	// before to run check context
	select {
	case <-ctx.Done():
		return
	default:
	}

	node, ok := reg.Get(start.Node)
	if !ok {
		log.Ctx(ctx).Error().Msgf("node %s not found", start.Node)

		return
	}

	// is disabled node
	if node.IsDisabled() {
		return
	}

	// add nodeID to log context
	ctx = log.Ctx(ctx).With().Str("nodeID", node.NodeID()).Logger().WithContext(ctx)
	// log debug
	log.Ctx(ctx).Debug().Msgf("running [%s]", node.GetType())

	outputDatas, err := node.Run(ctx, &reg.wgx, reg.appStore, value, start.Output)
	if err != nil {
		if errors.Is(err, ErrStopGoroutine) {
			return
		}

		log.Ctx(ctx).Error().Err(err).Msgf("%v cannot run", node.GetType())

		reg.AddError(fmt.Errorf("%s cannot run; nodeID=[%s]: %w", node.GetType(), node.NodeID(), err))

		return
	}

	log.Ctx(ctx).Debug().Msgf("complete [%s]", node.GetType())

	// direct go to output
	if outputDatasRespond, ok := outputDatas.(NodeDirectGo); ok {
		branch(ctx, node.Next(0), reg, outputDatasRespond.IsDirectGo())

		return
	}

	// check there is a respond interface
	if outputDatasRespond, ok := outputDatas.(NodeRetRespond); ok {
		// only one respond protection
		reg.mutex.Lock()
		if reg.respondChanActive {
			reg.respondChanActive = false

			if reg.respondChan != nil {
				reg.respondChan <- outputDatasRespond.GetRespond()
			}
		}
		reg.mutex.Unlock()

		return
	}

	// returning more than one data
	// call everything as for loop
	if outputDatasFor, ok := outputDatas.(NodeRetDatas); ok {
		datas := outputDatasFor.GetBinaryDatas()
		for i := range datas {
			branch(ctx, node.Next(0), reg, &nodeRetOutput{datas[i]})
		}

		return
	}

	// selection list for output
	if outputDatasSelection, ok := outputDatas.(NodeRetSelection); ok {
		for _, i := range outputDatasSelection.GetSelection() {
			branch(ctx, node.Next(i), reg, outputDatas)
		}

		return
	}

	// just one output group
	branch(ctx, node.Next(0), reg, outputDatas)
}
