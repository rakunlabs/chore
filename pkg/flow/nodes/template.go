package nodes

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var templateType = "template"

// Template node has one input and one output.
type Template struct {
	templateName string
	typeName     string
	inputs       []flow.Inputs
	inputHolder  map[string]interface{}
	outputs      []string
	content      []byte
	wait         int
	fetched      bool
	lock         sync.Mutex
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Template) Run(_ context.Context, reg *registry.AppStore, value []byte) ([]byte, error) {
	n.lock.Lock()
	n.wait--

	if n.wait < 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	// this block should be stay here in locking area
	// save value if wait is zero and more
	if err := yaml.Unmarshal(value, &n.inputHolder); err != nil {
		n.lock.Unlock()

		return nil, fmt.Errorf("template cannot unmarhal: %v", err)
	}

	if n.wait != 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	n.lock.Unlock()

	payload, err := reg.Template.Ext(n.inputHolder, string(n.content))
	if err != nil {
		err = fmt.Errorf("template cannot render: %v", err)
	}

	return payload, err
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

func (n *Template) Validate() error {
	return nil
}

func (n *Template) Next() []string {
	return n.outputs
}

func (n *Template) ActiveInput(node string) {
	for i := range n.inputs {
		if n.inputs[i].Node == node {
			n.inputs[i].Active = true
			n.wait++

			break
		}
	}
}

func (n *Template) CheckData() string {
	return ""
}

func NewTemplate(data flow.NodeData) flow.Noder {
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
