package nodes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var scriptType = "script"

type inputHolderS struct {
	input string
	value []byte
}

type ScriptRet struct {
	selection []int
	output    []byte
}

func (r *ScriptRet) GetBinaryData() []byte {
	return r.output
}

func (r *ScriptRet) GetSelection() []int {
	return r.selection
}

var _ flow.NodeRetSelection = &ScriptRet{}

// Script node has many input and one output.
type Script struct {
	script      string
	typeName    string
	inputs      []flow.Inputs
	inputHolder []inputHolderS
	outputs     [][]flow.Connection
	wait        int
	lock        sync.Mutex
	checked     bool
	directRun   bool
}

// selection 0 is false.
func (n *Script) Run(ctx context.Context, reg *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
	var inputValues []inputHolderS

	if !n.directRun {
		n.lock.Lock()
		n.wait--

		if n.wait < 0 {
			n.lock.Unlock()

			return nil, flow.ErrStopGoroutine
		}

		// this block should be stay here in locking area
		// save value if wait is zero and more
		n.inputHolder = append(n.inputHolder, inputHolderS{value: value.GetBinaryData(), input: input})

		if n.wait != 0 {
			n.lock.Unlock()

			return nil, flow.ErrStopGoroutine
		}

		n.lock.Unlock()
		inputValues = n.inputHolder
	} else {
		inputValues = []inputHolderS{{value: value.GetBinaryData(), input: input}}
	}

	scriptRunner := goja.New()

	// set script special functions
	if err := setScriptFuncs(scriptRunner); err != nil {
		return nil, err
	}

	if _, err := scriptRunner.RunString(n.script); err != nil {
		return nil, fmt.Errorf("script cannot read: %v", err)
	}

	mainScript, ok := goja.AssertFunction(scriptRunner.Get("main"))
	if !ok {
		return nil, fmt.Errorf("main function not found")
	}

	// sort inputholder by input name, it effects to function arguments order
	sort.Slice(inputValues, func(i, j int) bool {
		return inputValues[i].input < inputValues[j].input
	})

	passValues := []goja.Value{}
	for i := range inputValues {
		passValues = append(passValues, scriptRunner.ToValue(inputValues[i].value))
	}

	var retVal interface{}

	returnToFalse := false

	res, err := mainScript(goja.Undefined(), passValues...)
	if err != nil {
		var jserr *goja.Exception

		if errors.As(err, &jserr) {
			retVal = jserr.Value().Export()
			returnToFalse = true
		} else {
			return nil, fmt.Errorf("main function run: %v", err)
		}
	} else {
		retVal = res.Export()
	}

	var returnRes []byte

	if retVal != nil {
		switch exportValTyped := retVal.(type) {
		case map[string]interface{}, []interface{}:
			exportValM, err := json.Marshal(exportValTyped)
			if err != nil {
				return nil, fmt.Errorf("cannot marshal exported value: %v", err)
			}

			returnRes = exportValM
		case []byte:
			returnRes = exportValTyped
		default:
			returnRes = []byte(fmt.Sprint(exportValTyped))
		}
	}

	if returnToFalse {
		return &ScriptRet{
			selection: []int{0, 2},
			output:    returnRes,
		}, nil
	}

	return &ScriptRet{
		selection: []int{1, 2},
		output:    returnRes,
	}, nil
}

func (n *Script) GetType() string {
	return n.typeName
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

func (n *Script) Validate() error {
	return nil
}

func (n *Script) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Script) NextCount() int {
	return len(n.outputs)
}

func (n *Script) ActiveInput(node string) {
	for i := range n.inputs {
		if n.inputs[i].Node == node {
			n.inputs[i].Active = true
			n.wait++
			n.directRun = n.wait == 1

			break
		}
	}
}

func (n *Script) CheckData() string {
	return ""
}

func (n *Script) Check() {
	n.checked = true
}

func (n *Script) IsChecked() bool {
	return n.checked
}

func NewScript(_ context.Context, data flow.NodeData) flow.Noder {
	inputs := flow.PrepareInputs(data.Inputs)

	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	script, _ := data.Data["script"].(string)

	return &Script{
		typeName: scriptType,
		inputs:   inputs,
		outputs:  outputs,
		script:   script,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[scriptType] = NewScript
}

// ///////////////////////////////////

func toObject(v []byte) interface{} {
	var m interface{}

	_ = yaml.Unmarshal(v, &m)

	return m
}

func toString(v []byte) string {
	return string(v)
}

func sleep(length string) {
	duration, err := time.ParseDuration(length)
	if err != nil {
		panic(err)
	}

	time.Sleep(duration)
}

type commands struct {
	fn   interface{}
	name string
}

var commandList = []commands{
	{
		fn:   toObject,
		name: "toObject",
	},
	{
		fn:   toString,
		name: "toString",
	},
	{
		fn:   sleep,
		name: "sleep",
	},
}

func setScriptFuncs(runner *goja.Runtime) error {
	for _, v := range commandList {
		// custom functions set
		if err := runner.Set(v.name, v.fn); err != nil {
			return fmt.Errorf("%s command cannot set: %w", v.name, err)
		}
	}

	return nil
}
