package nodes

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gorm.io/gorm"
)

var controlType = "control"

type ControlRet struct {
	respond flow.Respond
}

func (r *ControlRet) GetBinaryData() []byte {
	return r.respond.Data
}

func (r *ControlRet) GetRespond() flow.Respond {
	return r.respond
}

var _ flow.NodeRetRespond = &ControlRet{}

// Control node has one input and one output.
type Control struct {
	controlName  string
	endpointName string
	inputs       []flow.Inputs
	outputs      [][]flow.Connection
	checked      bool
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Control) Run(ctx context.Context, reg *registry.AppStore, value flow.NodeRet, _ string) (flow.NodeRet, error) {
	control := models.Control{}

	query := reg.DB.WithContext(ctx).Where("name = ?", n.controlName)
	result := query.First(&control)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("control not found %s; %w", n.controlName, result.Error)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch %s; %w", n.controlName, result.Error)
	}

	content, err := base64.StdEncoding.DecodeString(control.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode content; %w", err)
	}

	log.Ctx(ctx).Info().Msgf("internal call control=[%s] endpoint=[%s]", control.Name, n.endpointName)
	nodesReg, err := flow.StartFlow(ctx, control.Name, n.endpointName, content, reg, value.GetBinaryData())
	if errors.Is(err, flow.ErrEndpointNotFound) {
		return nil, fmt.Errorf("endpoint not found %s; %w", n.endpointName, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to start [control:%s;endpoint:%s] content; %w", n.controlName, n.endpointName, err)
	}

	respondChan := nodesReg.GetChan()
	if respondChan == nil {
		return nil, nil
	}

	valueChan := <-respondChan

	return &ControlRet{valueChan}, nil
}

func (n *Control) Special(_ interface{}) interface{} {
	return nil
}

func (n *Control) GetType() string {
	return controlType
}

func (n *Control) Fetch(ctx context.Context, db *gorm.DB) error {
	return nil
}

func (n *Control) IsFetched() bool {
	return true
}

func (n *Control) IsRespond() bool {
	return false
}

func (n *Control) Validate() error {
	return nil
}

func (n *Control) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Control) NextCount() int {
	return len(n.outputs)
}

func (n *Control) ActiveInput(string) {}

func (n *Control) CheckData() string {
	return ""
}

func (n *Control) Check() {
	n.checked = true
}

func (n *Control) IsChecked() bool {
	return n.checked
}

func NewControl(_ context.Context, data flow.NodeData) flow.Noder {
	inputs := flow.PrepareInputs(data.Inputs)
	outputs := flow.PrepareOutputs(data.Outputs)

	controlName, _ := data.Data["control"].(string)
	endpointName, _ := data.Data["endpoint"].(string)

	return &Control{
		inputs:       inputs,
		outputs:      outputs,
		controlName:  controlName,
		endpointName: endpointName,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[controlType] = NewControl
}
