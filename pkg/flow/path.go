package flow

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	ErrStopGoroutine    = errors.New("stop goroutine")
	ErrEndpointNotFound = errors.New("endpoint not found")
)

func VisitAndFetch(ctx context.Context, reg *NodesReg) error {
	starts := reg.GetStarts()

	if starts == nil {
		return fmt.Errorf("%s %w", reg.startName, ErrEndpointNotFound)
	}

	return validateFetch(ctx, "", starts, reg)
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

		// activate this input
		node.ActiveInput(current)

		if node.IsChecked() {
			continue
		}

		// fetch and validation

		if err := node.Validate(); err != nil {
			return fmt.Errorf("ID %s, %s validate failed: %v", output, node.GetType(), err)
		}

		if err := node.Fetch(ctx, reg.appStore.DB); err != nil {
			return fmt.Errorf("ID %s, %s fetch failed: %v", output, node.GetType(), err)
		}

		// respond channel activate
		if node.IsRespond() {
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
	var stuckCtxCancel context.CancelFunc
	reg.stuckCtx, stuckCtxCancel = context.WithCancel(ctx)
	reg.stuckChan = make(chan bool, 1)

	reg.wgx.Add(1)
	go func() {
		defer reg.wgx.Done()

		for {
			select {
			case check := <-reg.stuckChan:
				if check {
					// all job is finished
					stuckCtxCancel()

					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// change waitgroup to check all job is finished
	branch(ctx, starts, reg, &nodeRetOutput{firstValue})

	reg.wgx.Wait()

	if reg.respondChan != nil {
		if reg.respondChanActive {
			// send error to channel
			buff := bytes.Buffer{}
			for _, err := range reg.errors {
				buff.WriteString("[")
				buff.WriteString(err.Error())
				buff.WriteString("]")
			}

			reg.respondChan <- Respond{
				Data:    buff.Bytes(),
				IsError: true,
			}
		}

		close(reg.respondChan)
	}
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
	defer reg.wgx.Done()
	defer reg.UpdateStuck(CountTotalDecrease, true)

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
