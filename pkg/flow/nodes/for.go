package nodes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/worldline-go/chore/pkg/flow"
	"github.com/worldline-go/chore/pkg/registry"
)

var forLoopType = "forLoop"

// ForLoop node has one input and one output.
// Not need to wait other inputs.
type ForLoop struct {
	expression string
	outputs    [][]flow.Connection
	checked    bool
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

func (n *ForLoop) Run(_ context.Context, _ *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
	scriptRunner := goja.New()

	var m interface{}
	if value.GetBinaryData() != nil {
		if err := yaml.Unmarshal(value.GetBinaryData(), &m); err != nil {
			m = value.GetBinaryData()
		}
	}

	if err := scriptRunner.Set("data", m); err != nil {
		return nil, fmt.Errorf("cannot set data in script: %v", err)
	}

	gojaV, err := scriptRunner.RunString(n.expression)
	if err != nil {
		return nil, fmt.Errorf("cannot run loop value: %v", err)
	}

	var v [][]byte

	if _, ok := gojaV.Export().([]interface{}); ok {
		for _, exportVal := range gojaV.Export().([]interface{}) {
			switch exportValTyped := exportVal.(type) {
			case map[string]interface{}, []interface{}:
				exportValM, err := json.Marshal(exportValTyped)
				if err != nil {
					return nil, fmt.Errorf("cannot marshal exported value: %v", err)
				}

				v = append(v, exportValM)
			case []byte:
				v = append(v, exportValTyped)
			default:
				v = append(v, []byte(fmt.Sprint(exportValTyped)))
			}
		}
	}

	return &ForRet{output: v}, nil
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

func NewForLoop(_ context.Context, data flow.NodeData) flow.Noder {
	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	expression, _ := data.Data["for"].(string)

	return &ForLoop{
		outputs:    outputs,
		expression: expression,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[forLoopType] = NewForLoop
}
