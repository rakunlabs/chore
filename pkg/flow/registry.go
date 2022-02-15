package flow

import (
	"context"
	"sync"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

type Respond struct {
	Data    []byte `json:"data"`
	Status  int    `json:"status"`
	IsError bool   `json:"isError"`
}

// NodesReg hold concreate information of nodes and start points.
type NodesReg struct {
	ctx               context.Context
	appStore          *registry.AppStore
	reg               map[string]Noder
	respondChan       chan Respond
	controlName       string
	startName         string
	starts            []string
	mutex             sync.RWMutex
	wg                sync.WaitGroup
	wgx               sync.WaitGroup
	respondChanActive bool
}

func NewNodesReg(ctx context.Context, controlName, startName string, appStore *registry.AppStore) *NodesReg {
	return &NodesReg{
		controlName: controlName,
		startName:   startName,
		reg:         make(map[string]Noder),
		appStore:    appStore,
		ctx:         ctx,
		respondChan: make(chan Respond, 1),
	}
}

func (r *NodesReg) GetChan() <-chan Respond {
	if r.respondChanActive {
		return r.respondChan
	}

	return nil
}

func (r *NodesReg) Wait() {
	r.wg.Wait()
}

func (r *NodesReg) Get(number string) Noder {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.reg[number]
}

func (r *NodesReg) GetStarts() []string {
	return r.starts
}

// Set a concreate node to registry.
// Number is a node number like 2, 4.
func (r *NodesReg) Set(number string, node Noder) {
	// checkdata usable for starter nodes like endpoint
	if node.CheckData() == r.startName {
		r.starts = append(r.starts, number)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.reg[number] = node
}
