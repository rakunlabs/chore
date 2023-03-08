package nodes

import (
	"context"
	"sort"
	"sync"

	"github.com/worldline-go/chore/pkg/flow"
	"github.com/worldline-go/chore/pkg/flow/convert"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/script/js"
	"github.com/worldline-go/chore/pkg/transfer"

	"gorm.io/gorm"
)

var scriptType = "script"

type inputHolderS struct {
	value interface{}
	input string
}

type ScriptRet struct {
	selection    []int
	output       []byte
	outputValues []byte
}

func (r *ScriptRet) GetBinaryData() []byte {
	return r.output
}

func (r *ScriptRet) GetSelection() []int {
	return r.selection
}

func (r *ScriptRet) GetBinaryValues() []byte {
	return r.outputValues
}

var (
	_ flow.NodeRetSelection = (*ScriptRet)(nil)
	_ flow.NodeRetValues    = (*ScriptRet)(nil)
)

// Script node has many input and one output.
type Script struct {
	script       string
	inputs       []flow.Inputs
	inputCounter map[string]struct{}
	inputHolder  map[string]inputHolderS
	outputs      [][]flow.Connection
	lock         sync.Mutex
	checked      bool
	disabled     bool
	nodeID       string
	tags         []string
}

// selection 0 is false.
//
//nolint:lll // false positive
func (n *Script) Run(ctx context.Context, _ *sync.WaitGroup, reg *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
	var transferValue interface{}
	if value.GetBinaryData() != nil {
		transferValue = transfer.BytesToData(value.GetBinaryData())
	}

	var inputValues []inputHolderS

	if len(n.inputCounter) > 1 {
		n.lock.Lock()

		// replace input value, multiple call on same input will replace value
		n.inputHolder[input] = inputHolderS{value: transferValue, input: input}

		if len(n.inputCounter) > len(n.inputHolder) {
			n.lock.Unlock()

			return nil, flow.ErrStopGoroutine
		}

		inputValues = make([]inputHolderS, 0, len(n.inputHolder))
		for _, v := range n.inputHolder {
			inputValues = append(inputValues, v)
		}

		n.lock.Unlock()
	} else {
		inputValues = []inputHolderS{{value: transferValue, input: input}}
	}

	// sort inputholder by input name, it effects to function arguments order
	sort.Slice(inputValues, func(i, j int) bool {
		return inputValues[i].input < inputValues[j].input
	})

	// create script runner
	runner := js.NewGoja()

	var valueToPass interface{}
	// value for some nodes
	setValue := func(v interface{}) {
		valueToPass = v
	}

	runner.SetFunction("setValue", setValue)

	inputValuesInterface := make([]interface{}, len(inputValues))
	for i, v := range inputValues {
		inputValuesInterface[i] = v.value
	}

	result, err := runner.RunScript(ctx, n.script, inputValuesInterface)
	if err != nil {
		return &ScriptRet{ //nolint:nilerr // different kind of error
			selection: []int{0, 2},
			output:    result,
		}, nil
	}

	return &ScriptRet{
		selection:    []int{1, 2},
		output:       result,
		outputValues: transfer.DataToBytes(valueToPass),
	}, nil
}

func (n *Script) GetType() string {
	return scriptType
}

func (n *Script) Fetch(_ context.Context, _ *gorm.DB) error {
	return nil
}

func (n *Script) IsFetched() bool {
	return true
}

func (n *Script) IsRespond() bool {
	return false
}

func (n *Script) Validate(_ context.Context) error {
	return nil
}

func (n *Script) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Script) NextCount() int {
	return len(n.outputs)
}

func (n *Script) IsDisabled() bool {
	return n.disabled
}

func (n *Script) ActiveInput(node string, tags map[string]struct{}) {
	if !convert.IsTagsEnabled(n.tags, tags) {
		n.disabled = true

		return
	}

	for i := range n.inputs {
		if n.inputs[i].Node == node {
			if !n.inputs[i].Active {
				n.inputs[i].Active = true
				n.inputCounter[n.inputs[i].InputName] = struct{}{}
			}
		}
	}
}

func (n *Script) Check() {
	n.checked = true
}

func (n *Script) IsChecked() bool {
	return n.checked
}

func (n *Script) NodeID() string {
	return n.nodeID
}

func (n *Script) Tags() []string {
	return n.tags
}

func NewScript(_ context.Context, _ *flow.NodesReg, data flow.NodeData, nodeID string) (flow.Noder, error) {
	inputs := flow.PrepareInputs(data.Inputs)

	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	script, _ := data.Data["script"].(string)
	tags := convert.GetList(data.Data["tags"])

	return &Script{
		inputs:       inputs,
		outputs:      outputs,
		script:       script,
		nodeID:       nodeID,
		inputCounter: make(map[string]struct{}),
		inputHolder:  make(map[string]inputHolderS),
		tags:         tags,
	}, nil
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[scriptType] = NewScript
}
