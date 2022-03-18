package nodes

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"github.com/dop251/goja"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var scriptType = "script"

type inputHolderS struct {
	input string
	value []byte
}

// Script node has many input and one output.
type Script struct {
	script      string
	typeName    string
	inputs      []flow.Inputs
	inputHolder []inputHolderS
	outputs     []flow.Connection
	wait        int
	lock        sync.Mutex
}

func tomap(v []byte) map[string]interface{} {
	var m map[string]interface{}

	_ = yaml.Unmarshal(v, &m)

	return m
}

func tostring(v []byte) string {
	return string(v)
}

func (n *Script) Run(_ context.Context, reg *registry.AppStore, value []byte, input string) ([]byte, error) {
	n.lock.Lock()
	n.wait--

	if n.wait < 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	// this block should be stay here in locking area
	// save value if wait is zero and more
	n.inputHolder = append(n.inputHolder, inputHolderS{value: value, input: input})

	if n.wait != 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	n.lock.Unlock()

	scriptRunner := goja.New()

	// custom functions set
	if err := scriptRunner.Set("tomap", tomap); err != nil {
		return nil, fmt.Errorf("tomap command cannot set: %v", err)
	}

	if err := scriptRunner.Set("tostring", tostring); err != nil {
		return nil, fmt.Errorf("tostring command cannot set: %v", err)
	}

	// end custom functions set

	if _, err := scriptRunner.RunString(n.script); err != nil {
		return nil, fmt.Errorf("script cannot read: %v", err)
	}

	mainScript, ok := goja.AssertFunction(scriptRunner.Get("main"))
	if !ok {
		return nil, fmt.Errorf("main function not found")
	}

	// sort inputholder by input name, it effects to function arguments order
	sort.Slice(n.inputHolder, func(i, j int) bool {
		return n.inputHolder[i].input < n.inputHolder[j].input
	})

	log.Debug().Msgf("%#v", n.inputHolder)

	passValues := []goja.Value{}
	for i := range n.inputHolder {
		passValues = append(passValues, scriptRunner.ToValue(n.inputHolder[i].value))
	}

	res, err := mainScript(goja.Undefined(), passValues...)
	if err != nil {
		return nil, fmt.Errorf("main function run: %v", err)
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

func (n *Script) Next() []flow.Connection {
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

	outputs := data.Outputs["output_1"].Connections

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
