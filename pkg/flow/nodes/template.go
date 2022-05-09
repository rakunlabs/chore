package nodes

import (
	"context"
	"encoding/base64"
	"fmt"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gopkg.in/yaml.v3"
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
	typeName     string
	inputs       []flow.Inputs
	outputs      [][]flow.Connection
	content      []byte
	fetched      bool
	checked      bool
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Template) Run(_ context.Context, reg *registry.AppStore, value flow.NodeRet, _ string) (flow.NodeRet, error) {
	var v map[string]interface{}
	if err := yaml.Unmarshal(value.GetBinaryData(), &v); err != nil {
		return nil, fmt.Errorf("template cannot unmarhal: %v", err)
	}

	payload, err := reg.Template.Ext(v, string(n.content))
	if err != nil {
		err = fmt.Errorf("template cannot render: %v", err)
	}

	return &TemplateRet{payload}, err
}

func (n *Template) Special(_ interface{}) interface{} {
	return nil
}

func (n *Template) GetType() string {
	return n.typeName
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

func (n *Template) CheckData() string {
	return ""
}

func (n *Template) Check() {
	n.checked = true
}

func (n *Template) IsChecked() bool {
	return n.checked
}

func NewTemplate(_ context.Context, data flow.NodeData) flow.Noder {
	inputs := flow.PrepareInputs(data.Inputs)

	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	templateName, _ := data.Data["template"].(string)

	return &Template{
		typeName:     templateType,
		inputs:       inputs,
		outputs:      outputs,
		templateName: templateName,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[templateType] = NewTemplate
}
