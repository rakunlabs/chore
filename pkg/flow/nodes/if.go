package nodes

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

var ifCaseType = "ifCase"

// Ifcase node has one input and one output.
// Not need to wait other inputs.
type IfCase struct {
	typeName   string
	expression string
	outputs    [][]flow.Connection
}

func (n *IfCase) Run(ctx context.Context, _ *registry.AppStore, value []byte, input string) ([][]byte, error) {
	scriptRunner := goja.New()

	if err := scriptRunner.Set("data", toObject(value)); err != nil {
		return nil, fmt.Errorf("cannot set data in script: %v", err)
	}

	gojaV, err := scriptRunner.RunString(n.expression)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("cannot run loop value, passing as false: %v", err)

		return [][]byte{value, nil}, nil
	}

	if gojaV.ToBoolean() {
		return [][]byte{nil, value}, nil
	}

	return [][]byte{value, nil}, nil
}

func (n *IfCase) GetType() string {
	return n.typeName
}

func (n *IfCase) Fetch(_ context.Context, _ *gorm.DB) error {
	return nil
}

func (n *IfCase) IsFetched() bool {
	return true
}

func (n *IfCase) Validate() error {
	return nil
}

func (n *IfCase) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *IfCase) NextCount() int {
	return len(n.outputs)
}

func (n *IfCase) ActiveInput(node string) {}

func (n *IfCase) CheckData() string {
	return ""
}

func NewIfCase(data flow.NodeData) flow.Noder {
	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	expression, _ := data.Data["if"].(string)

	return &IfCase{
		typeName:   ifCaseType,
		outputs:    outputs,
		expression: expression,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[ifCaseType] = NewIfCase
}
