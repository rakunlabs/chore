package flow

import (
	"errors"
	"fmt"
	"net/http"

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
	for _, output := range outputs {
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

		if err := node.Validate(); err != nil {
			return fmt.Errorf("ID %s, %s validate failed: %v", output, node.GetType(), err)
		}

		if err := node.Fetch(reg.ctx, reg.appStore.DB); err != nil {
			return fmt.Errorf("ID %s, %s fetch failed: %v", output, node.GetType(), err)
		}

		// activate this input
		node.ActiveInput(current)

		// respond channel activate
		if node.GetType() == "respond" {
			reg.respondChanActive = true
		}

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
	branch(starts, reg, firstValue)

	reg.wgx.Wait()
	close(reg.respondChan)
}

// branch values already designed to same length of nexts.
func branch(nexts []Connection, reg *NodesReg, value []byte) {
	for _, next := range nexts {
		// going goroutine to prevent too much recursive call
		reg.wgx.Add(1)

		go branchRun(next, reg, value)
	}
}

func branchRun(start Connection, reg *NodesReg, value []byte) {
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

	// special cases

	switch node.GetType() {
	case "respond":
		// only one respond protection
		reg.mutex.Lock()
		if reg.respondChanActive {
			reg.respondChanActive = false

			reg.respondChan <- Respond{
				// respond has one output
				Data:   outputDatas[0],
				Status: http.StatusOK,
			}
		}
		reg.mutex.Unlock()

		return

	case "forLoop":
		for _, outputData := range outputDatas {
			branch(node.Next(0), reg, outputData)
		}

		return

	case "request", "script", "ifCase":
		// first data is error
		if outputDatas[0] != nil {
			branch(node.Next(0), reg, outputDatas[0])

			return
		}

		branch(node.Next(1), reg, outputDatas[1])

		return
	}

	// separate data for branchs
	// first data goes to first output (output_1)
	for i, outputData := range outputDatas {
		branch(node.Next(i), reg, outputData)
	}
}
