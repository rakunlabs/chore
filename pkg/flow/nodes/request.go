package nodes

import (
	"context"
	"fmt"
	"sync"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var requestType = "request"

type inputHolderRequest struct {
	value []byte
	data  []byte
}

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
	inputHolder   inputHolderRequest
	lock          sync.Mutex
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Request) Run(ctx context.Context, reg *registry.AppStore, value []byte, input string) ([]byte, error) {
	n.lock.Lock()
	n.wait--

	if n.wait < 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	// this block should be stay here in locking area
	// save value if wait is zero and more
	// input_1 is value
	if input == "input_1" {
		n.inputHolder.value = value
	} else {
		n.inputHolder.data = value
	}

	if n.wait != 0 {
		n.lock.Unlock()

		return nil, flow.ErrStopGoroutine
	}

	n.lock.Unlock()

	// check value and render it
	if n.inputHolder.value != nil {
		var requestValues map[string]interface{}
		if err := yaml.Unmarshal(n.inputHolder.value, &requestValues); err != nil {
			return nil, err
		}

		// render url
		payload, err := reg.Template.Ext(requestValues, n.url)
		if err != nil {
			return nil, fmt.Errorf("template cannot render: %v", err)
		}

		n.url = string(payload)

		// render method
		payload, err = reg.Template.Ext(requestValues, n.method)
		if err != nil {
			return nil, fmt.Errorf("template cannot render: %v", err)
		}

		n.method = string(payload)

		// render headers
		payload, err = reg.Template.Ext(requestValues, n.addHeadersRaw)
		if err != nil {
			return nil, fmt.Errorf("template cannot render: %v", err)
		}

		n.addHeadersRaw = string(payload)
	}

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
		n.inputHolder.data,
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
