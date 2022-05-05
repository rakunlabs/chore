package flow

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

var (
	ErrStopGoroutine    = errors.New("stop goroutine")
	ErrEndpointNotFound = errors.New("endpoint not found")
)

func VisitAndFetch(reg *NodesReg) error {
	starts := reg.GetStarts()

	if starts == nil {
		return fmt.Errorf("%s %w", reg.startName, ErrEndpointNotFound)
	}

	return validateFetch("", starts, reg)
}

func validateFetch(current string, outputs []Connection, reg *NodesReg) error {
	// log.Debug().Msgf("current %s", current)
	for _, output := range outputs {
		// log.Debug().Msgf("validating [%v]", output)
		// before to run check context
		select {
		case <-reg.ctx.Done():
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

		if err := node.Fetch(reg.ctx, reg.appStore.DB); err != nil {
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
			if err := validateFetch(output.Node, node.Next(i), reg); err != nil {
				return err
			}
		}
	}

	return nil
}

func GoAndRun(reg *NodesReg, firstValue []byte) {
	defer reg.wg.Done()
	starts := reg.GetStarts()

	// change waitgroup to check all job is finished
	branch(starts, reg, &nodeRetOutput{firstValue})

	reg.wgx.Wait()
	close(reg.respondChan)
}

// branch values already designed to same length of nexts.
func branch(nexts []Connection, reg *NodesReg, value NodeRet) {
	for _, next := range nexts {
		// going goroutine to prevent too much recursive call
		reg.wgx.Add(1)

		go branchRun(next, reg, value)
	}
}

func branchRun(start Connection, reg *NodesReg, value NodeRet) {
	defer reg.wgx.Done()

	// before to run check context
	select {
	case <-reg.ctx.Done():
		return
	default:
	}

	node, ok := reg.Get(start.Node)
	if !ok {
		log.Ctx(reg.ctx).Error().Msgf("node %s not found", start.Node)

		return
	}

	outputDatas, err := node.Run(reg.ctx, reg.appStore, value, start.Output)
	if err != nil {
		if errors.Is(err, ErrStopGoroutine) {
			return
		}

		log.Ctx(reg.ctx).Error().Err(err).Msgf("%v cannot run", node.GetType())

		return
	}

	// check there is a repond interface
	if _, ok := outputDatas.(NodeRetRespond); ok {
		// only one respond protection
		reg.mutex.Lock()
		if reg.respondChanActive {
			reg.respondChanActive = false

			reg.respondChan <- outputDatas.(NodeRetRespond).GetRespond()
		}
		reg.mutex.Unlock()

		return
	}

	// returning more than one data
	// call everything as for loop
	if _, ok := outputDatas.(NodeRetDatas); ok {
		for _, outputData := range outputDatas.(NodeRetDatas).GetBinaryDatas() {
			branch(node.Next(0), reg, &nodeRetOutput{outputData})
		}

		return
	}

	// selection list for output
	if _, ok := outputDatas.(NodeRetSelection); ok {
		for _, i := range outputDatas.(NodeRetSelection).GetSelection() {
			branch(node.Next(i), reg, outputDatas)
		}

		return
	}

	// just one output
	branch(node.Next(0), reg, outputDatas)
}
