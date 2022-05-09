package nodes

import (
	"context"
	"fmt"

	"github.com/dop251/goja"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

var ifCaseType = "ifCase"

type IfRet struct {
	output    []byte
	selection []int
}

func (r *IfRet) GetBinaryData() []byte {
	return r.output
}

func (r *IfRet) GetSelection() []int {
	return r.selection
}

var _ flow.NodeRetSelection = &IfRet{}

// Ifcase node has one input and one output.
// Not need to wait other inputs.
type IfCase struct {
	typeName   string
	expression string
	outputs    [][]flow.Connection
	checked    bool
}

// selection 0 is false.
func (n *IfCase) Run(ctx context.Context, _ *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
	scriptRunner := goja.New()

	var m interface{}
	if value.GetBinaryData() != nil {
		if err := yaml.Unmarshal(value.GetBinaryData(), &m); err != nil {
			m = value.GetBinaryData()
		}
	}

	// set script special functions
	if err := setScriptFuncs(scriptRunner); err != nil {
		return nil, err
	}

	if err := scriptRunner.Set("data", m); err != nil {
		return nil, fmt.Errorf("cannot set data in script: %v", err)
	}

	gojaV, err := scriptRunner.RunString(n.expression)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("cannot run loop value, passing as false: %v", err)

		return &IfRet{
			output:    value.GetBinaryData(),
			selection: []int{0},
		}, nil
	}

	if gojaV.ToBoolean() {
		return &IfRet{
			output:    value.GetBinaryData(),
			selection: []int{1},
		}, nil
	}

	return &IfRet{
		output:    value.GetBinaryData(),
		selection: []int{0},
	}, nil
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

func (n *IfCase) IsRespond() bool {
	return false
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

func (n *IfCase) ActiveInput(string) {}

func (n *IfCase) CheckData() string {
	return ""
}

func (n *IfCase) Check() {
	n.checked = true
}

func (n *IfCase) IsChecked() bool {
	return n.checked
}

func NewIfCase(_ context.Context, data flow.NodeData) flow.Noder {
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
