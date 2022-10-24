package nodes

import (
	"context"
	"fmt"
	"sync"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/chore/pkg/flow"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/script/js"
	"github.com/worldline-go/chore/pkg/transfer"
)

var forLoopType = "forLoop"

// ForLoop node has one input and one output.
// Not need to wait other inputs.
type ForLoop struct {
	expression string
	outputs    [][]flow.Connection
	checked    bool
	nodeID     string
}

type ForRet struct {
	output [][]byte
}

func (r *ForRet) GetBinaryData() []byte {
	return nil
}

func (r *ForRet) GetBinaryDatas() [][]byte {
	return r.output
}

var _ flow.NodeRetDatas = &ForRet{}

func (n *ForLoop) Run(ctx context.Context, _ *sync.WaitGroup, _ *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
	transferValue := transfer.BytesToData(value.GetBinaryData())

	runner := js.NewGoja()

	if err := runner.SetData(transferValue); err != nil {
		return nil, fmt.Errorf("cannot set data in script: %v", err)
	}

	gojaV, err := runner.RunString(n.expression)
	if err != nil {
		return nil, fmt.Errorf("cannot run loop value: %v", err)
	}

	var forValues [][]byte

	vExported := gojaV.Export()

	if vSlice, ok := vExported.([]interface{}); ok {
		for _, exportVal := range vSlice {
			forValues = append(forValues, transfer.DataToBytes(exportVal))
		}
	} else {
		log.Ctx(ctx).Warn().Msgf("for loop value is not a slice: %v", vExported)
	}

	if len(forValues) == 0 {
		return nil, flow.ErrStopGoroutine
	}

	return &ForRet{output: forValues}, nil
}

func (n *ForLoop) GetType() string {
	return forLoopType
}

func (n *ForLoop) Fetch(_ context.Context, _ *gorm.DB) error {
	return nil
}

func (n *ForLoop) IsFetched() bool {
	return true
}

func (n *ForLoop) IsRespond() bool {
	return false
}

func (n *ForLoop) Validate() error {
	return nil
}

func (n *ForLoop) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *ForLoop) NextCount() int {
	return len(n.outputs)
}

func (n *ForLoop) ActiveInput(string) {}

func (n *ForLoop) Check() {
	n.checked = true
}

func (n *ForLoop) IsChecked() bool {
	return n.checked
}

func (n *ForLoop) NodeID() string {
	return n.nodeID
}

func NewForLoop(_ context.Context, _ *flow.NodesReg, data flow.NodeData, nodeID string) (flow.Noder, error) {
	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	expression, _ := data.Data["for"].(string)

	return &ForLoop{
		outputs:    outputs,
		expression: expression,
		nodeID:     nodeID,
	}, nil
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[forLoopType] = NewForLoop
}
