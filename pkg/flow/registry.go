package flow

import (
	"context"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/worldline-go/chore/pkg/registry"
)

type Respond struct {
	Header map[string]interface{} `json:"header"`
	Data   []byte                 `json:"data"`
	Status int                    `json:"status"`
}

// NodesReg hold concreate information of nodes and start points.
type NodesReg struct {
	ctx               context.Context
	appStore          *registry.AppStore
	reg               map[string]Noder
	respondChan       chan Respond
	controlName       string
	startName         string
	method            string
	starts            []Connection
	mutex             sync.RWMutex
	wg                sync.WaitGroup
	wgx               sync.WaitGroup
	respondChanActive bool
}

func NewNodesReg(ctx context.Context, controlName, startName, method string, appStore *registry.AppStore) *NodesReg {
	// set new logger for reg and set it in ctx
	logReg := log.Ctx(ctx).With().Str("control", controlName).Str("endpoint", startName).Logger()

	return &NodesReg{
		controlName: controlName,
		startName:   startName,
		method:      method,
		reg:         make(map[string]Noder),
		appStore:    appStore,
		ctx:         logReg.WithContext(ctx),
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

func (r *NodesReg) Get(number string) (Noder, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	node, ok := r.reg[number]

	return node, ok
}

func (r *NodesReg) GetStarts() []Connection {
	return r.starts
}

// Set a concreate node to registry.
// Number is a node number like 2, 4.
func (r *NodesReg) Set(number string, node Noder) {
	// checkdata usable for starter nodes like endpoint
	if nodeEndpoint, ok := node.(NoderEndpoint); ok {
		if nodeEndpoint.Endpoint() == r.startName {
			for _, v := range nodeEndpoint.Methods() {
				if strings.ToUpper(v) == r.method {
					r.starts = append(r.starts, Connection{
						Node: number,
					})

					break
				}
			}
		}
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.reg[number] = node
}
