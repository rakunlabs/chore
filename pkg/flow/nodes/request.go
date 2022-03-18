package nodes

import (
	"context"
	"fmt"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var requestType = "request"

// Request node has one input and one output.
type Request struct {
	auth          string
	method        string
	url           string
	headers       map[string]interface{}
	addHeadersRaw string
	typeName      string
	inputs        []flow.Inputs
	outputs       []flow.Connection
	wait          int
	fetched       bool
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Request) Run(ctx context.Context, reg *registry.AppStore, value []byte, _ string) ([]byte, error) {
	var addHeaders map[string]interface{}
	if err := yaml.Unmarshal([]byte(n.addHeadersRaw), &addHeaders); err != nil {
		return nil, err
	}

	for k := range addHeaders {
		n.headers[k] = addHeaders[k]
	}

	response, err := reg.Client.Send(
		ctx,
		n.url,
		n.method,
		n.headers,
		value,
	)
	if err != nil {
		return nil, err
	}

	return response.Body, nil
}

func (n *Request) GetType() string {
	return n.typeName
}

func (n *Request) Fetch(ctx context.Context, db *gorm.DB) error {
	if n.auth == "" {
		n.fetched = true
		n.headers = make(map[string]interface{})

		return nil
	}

	getData := models.AuthPure{}

	query := db.WithContext(ctx).Model(&models.Auth{}).Where("name = ?", n.auth)
	result := query.First(&getData)

	if result.Error != nil {
		return fmt.Errorf("request fetch failed: %v", result.Error)
	}

	n.headers = getData.Headers

	n.fetched = true

	return nil
}

func (n *Request) IsFetched() bool {
	return n.fetched
}

func (n *Request) Validate() error {
	// correction
	if n.method == "" {
		n.method = "POST"
	}

	if n.url == "" {
		return fmt.Errorf("url is empty")
	}

	return nil
}

func (n *Request) Next() []flow.Connection {
	return n.outputs
}

func (n *Request) ActiveInput(node string) {
	for i := range n.inputs {
		if n.inputs[i].Node == node {
			n.inputs[i].Active = true
			n.wait++

			break
		}
	}
}

func (n *Request) CheckData() string {
	return ""
}

func NewRequest(data flow.NodeData) flow.Noder {
	inputs := make([]flow.Inputs, 0, len(data.Inputs))

	for _, input := range data.Inputs {
		for _, connection := range input.Connections {
			inputs = append(inputs, flow.Inputs{Node: connection.Node})
		}
	}

	outputs := data.Outputs["output_1"].Connections

	auth, _ := data.Data["auth"].(string)
	method, _ := data.Data["method"].(string)
	url, _ := data.Data["url"].(string)
	addHeadersRaw, _ := data.Data["headers"].(string)

	return &Request{
		typeName:      requestType,
		inputs:        inputs,
		outputs:       outputs,
		auth:          auth,
		method:        method,
		url:           url,
		addHeadersRaw: addHeadersRaw,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[requestType] = NewRequest
}
