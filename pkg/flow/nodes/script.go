package nodes

import (
	"context"
	"fmt"
	"sync"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var scriptType = "script"

// Script node has many input and one output.
type Script struct {
	script      string
	typeName    string
	inputs      []flow.Inputs
	inputHolder [][]byte
	outputs     []string
	wait        int
	lock        sync.Mutex
}

func Unmarshal(v []byte) map[string]interface{} {
	var m map[string]interface{}

	_ = yaml.Unmarshal(v, &m)

	return m
}

func (n *Script) Run(_ context.Context, reg *registry.AppStore, value []byte) ([]byte, error) {
	n.lock.Lock()
	n.wait--

	if n.wait < 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	// this block should be stay here in locking area
	// save value if wait is zero and more
	n.inputHolder = append(n.inputHolder, value)

	if n.wait != 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	n.lock.Unlock()

	scriptRunner := goja.New()

	if err := scriptRunner.Set("unmarshal", Unmarshal); err != nil {
		return nil, fmt.Errorf("unmarshal command cannot set: %v", err)
	}

	if _, err := scriptRunner.RunString(n.script); err != nil {
		return nil, fmt.Errorf("script cannot read: %v", err)
	}

	mainScript, ok := goja.AssertFunction(scriptRunner.Get("main"))
	if !ok {
		return nil, fmt.Errorf("parse function not found")
	}

	passValues := []goja.Value{}
	for i := range n.inputHolder {
		passValues = append(passValues, scriptRunner.ToValue(n.inputHolder[i]))
	}

	res, err := mainScript(goja.Undefined(), passValues...)
	if err != nil {
		return nil, fmt.Errorf("parse function run: %v", err)
	}

	var returnRes []byte
	switch v := res.Export().(type) {
	case []byte:
		returnRes = v
	default:
		returnRes = []byte(fmt.Sprintf("%v", res.Export()))
	}

	return returnRes, nil
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

func (n *Script) Validate() error {
	return nil
}

func (n *Script) Next() []string {
	return n.outputs
}

func (n *Script) ActiveInput(node string) {
	for i := range n.inputs {
		if n.inputs[i].Node == node {
			n.inputs[i].Active = true
			n.wait++

			break
		}
	}
}

func (n *Script) CheckData() string {
	return ""
}

func NewScript(data flow.NodeData) flow.Noder {
	inputs := make([]flow.Inputs, 0, len(data.Inputs))

	for _, input := range data.Inputs {
		for _, connection := range input.Connections {
			inputs = append(inputs, flow.Inputs{Node: connection.Node})
		}
	}

	outputs := make([]string, 0, len(data.Outputs["output_1"].Connections))
	for _, connection := range data.Outputs["output_1"].Connections {
		outputs = append(outputs, connection.Node)
	}

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
