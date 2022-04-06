package nodes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

var forLoopType = "forLoop"

// ForLoop node has one input and one output.
// Not need to wait other inputs.
type ForLoop struct {
	typeName   string
	expression string
	outputs    [][]flow.Connection
}

func (n *ForLoop) Run(_ context.Context, _ *registry.AppStore, value []byte, input string) ([][]byte, error) {
	scriptRunner := goja.New()

	if err := scriptRunner.Set("data", toObject(value)); err != nil {
		return nil, fmt.Errorf("cannot set data in script: %v", err)
	}

	gojaV, err := scriptRunner.RunString(n.expression)
	if err != nil {
		return nil, fmt.Errorf("cannot run loop value: %v", err)
	}

	var v [][]byte

	if _, ok := gojaV.Export().([]interface{}); ok {
		for _, exportVal := range gojaV.Export().([]interface{}) {
			gojaVByte, err := json.Marshal(exportVal)
			if err != nil {
				return nil, fmt.Errorf("cannot marshal exported value: %v", err)
			}

			v = append(v, gojaVByte)
		}
	}

	return v, nil
}

func (n *ForLoop) GetType() string {
	return n.typeName
}

func (n *ForLoop) Fetch(_ context.Context, _ *gorm.DB) error {
	return nil
}

func (n *ForLoop) IsFetched() bool {
	return true
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

func (n *ForLoop) ActiveInput(node string) {}

func (n *ForLoop) CheckData() string {
	return ""
}

func NewForLoop(data flow.NodeData) flow.Noder {
	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	expression, _ := data.Data["for"].(string)

	return &ForLoop{
		typeName:   forLoopType,
		outputs:    outputs,
		expression: expression,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[forLoopType] = NewForLoop
}
