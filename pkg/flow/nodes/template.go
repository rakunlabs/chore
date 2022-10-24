package nodes

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/worldline-go/chore/models"
	"github.com/worldline-go/chore/pkg/flow"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/transfer"

	"gorm.io/gorm"
)

var templateType = "template"

type TemplateRet struct {
	output []byte
}

func (r *TemplateRet) GetBinaryData() []byte {
	return r.output
}

// Template node has one input and one output.
type Template struct {
	templateName string
	inputs       []flow.Inputs
	outputs      [][]flow.Connection
	content      []byte
	fetched      bool
	checked      bool
	nodeID       string
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Template) Run(_ context.Context, _ *sync.WaitGroup, reg *registry.AppStore, value flow.NodeRet, _ string) (flow.NodeRet, error) {
	v := transfer.BytesToData(value.GetBinaryData())

	payload, err := reg.Template.ExecuteBytes(v, string(n.content))
	if err != nil {
		err = fmt.Errorf("template cannot render: %v", err)
	}

	return &TemplateRet{payload}, err
}

func (n *Template) Special(_ interface{}) interface{} {
	return nil
}

func (n *Template) GetType() string {
	return templateType
}

func (n *Template) Fetch(ctx context.Context, db *gorm.DB) error {
	if n.templateName == "" {
		return fmt.Errorf("template fetch failed: templateName empty")
	}

	getData := models.TemplatePure{}

	query := db.WithContext(ctx).Model(&models.Template{}).Where("name = ?", n.templateName)
	result := query.First(&getData)

	if result.Error != nil {
		return fmt.Errorf("template fetch failed: %v", result.Error)
	}

	content, err := base64.StdEncoding.DecodeString(getData.Content)
	if err != nil {
		return fmt.Errorf("template fetch failed: %v", err)
	}

	n.content = content

	n.fetched = true

	return nil
}

func (n *Template) IsFetched() bool {
	return n.fetched
}

func (n *Template) IsRespond() bool {
	return false
}

func (n *Template) Validate() error {
	return nil
}

func (n *Template) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Template) NextCount() int {
	return len(n.outputs)
}

func (n *Template) ActiveInput(string) {}

func (n *Template) Check() {
	n.checked = true
}

func (n *Template) IsChecked() bool {
	return n.checked
}

func (n *Template) NodeID() string {
	return n.nodeID
}

func NewTemplate(_ context.Context, _ *flow.NodesReg, data flow.NodeData, nodeID string) (flow.Noder, error) {
	inputs := flow.PrepareInputs(data.Inputs)

	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	templateName, _ := data.Data["template"].(string)

	return &Template{
		inputs:       inputs,
		outputs:      outputs,
		templateName: templateName,
		nodeID:       nodeID,
	}, nil
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[templateType] = NewTemplate
}
