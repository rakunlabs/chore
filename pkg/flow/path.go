package flow

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

var ErrStopGoroutine = errors.New("stop goroutine")

func VisitAndFetch(reg *NodesReg) error {
	starts := reg.GetStarts()

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

		node := reg.Get(output.Node)
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

		if err := validateFetch(output.Node, node.Next(), reg); err != nil {
			return err
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

	node := reg.Get(start.Node)

	outputData, err := node.Run(reg.ctx, reg.appStore, value, start.Output)
	if err != nil {
		if errors.Is(err, ErrStopGoroutine) {
			return
		}

		log.Error().Err(err).Msgf("%v cannot run", node.GetType())

		return
	}

	// only one respond protection
	if node.GetType() == "respond" {
		reg.mutex.Lock()
		if reg.respondChanActive {
			reg.respondChanActive = false

			reg.respondChan <- Respond{
				Data:   outputData,
				Status: http.StatusOK,
			}
		}
		reg.mutex.Unlock()
	}

	branch(node.Next(), reg, outputData)
}
