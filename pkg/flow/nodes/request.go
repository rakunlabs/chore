package nodes

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

var requestType = "request"

type inputHolderRequest struct {
	value []byte
	exist bool
}

type RequestRet struct {
	selection []int
	respond   flow.Respond
}

func (r *RequestRet) GetBinaryData() []byte {
	return r.respond.Data
}

func (r *RequestRet) GetSelection() []int {
	return r.selection
}

func (r *RequestRet) GetRespondData() flow.Respond {
	return r.respond
}

var (
	_ flow.NodeRetRespondData = &RequestRet{}
	_ flow.NodeRetSelection   = &RequestRet{}
)

// Request node has one input and one output.
type Request struct {
	lockCtx       context.Context
	lockCancel    context.CancelFunc
	headers       map[string]interface{}
	auth          string
	method        string
	url           string
	addHeadersRaw string
	typeName      string
	inputs        []flow.Inputs
	outputs       [][]flow.Connection
	inputHolder   inputHolderRequest
	mutex         sync.Mutex
	fetched       bool
	checked       bool
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Request) Run(ctx context.Context, reg *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
	// input_1 is value
	if input == flow.Input1 {
		// don't allow multiple inputs
		n.mutex.Lock()
		defer n.mutex.Unlock()

		if n.inputHolder.exist {
			return nil, flow.ErrStopGoroutine
		}

		n.inputHolder.value = value.GetBinaryData()
		n.inputHolder.exist = true

		// close context to allow to others continue process
		if n.lockCancel != nil {
			n.lockCancel()
		}

		return nil, flow.ErrStopGoroutine
	}

	if n.lockCtx != nil {
		select {
		case <-time.After(time.Hour * 1):
			log.Ctx(ctx).Warn().Msg("timeline exceded, terminated request")

			return nil, flow.ErrStopGoroutine
		case <-ctx.Done():
			log.Ctx(ctx).Warn().Msg("program closed, terminated request")

			return nil, flow.ErrStopGoroutine
		case <-n.lockCtx.Done():
		}
	}

	// check value and render it
	if n.inputHolder.value != nil {
		var requestValues map[string]interface{}
		if err := yaml.Unmarshal(n.inputHolder.value, &requestValues); err != nil {
			return nil, fmt.Errorf("failed to hnmarshal request values: %v", err)
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
		return nil, fmt.Errorf("faild unmarshal headers in request: %v", err)
	}

	for k := range addHeaders {
		n.headers[k] = addHeaders[k]
	}

	response, err := reg.Client.Send(
		ctx,
		n.url,
		n.method,
		n.headers,
		value.GetBinaryData(),
	)
	if err != nil {
		// return nil, fmt.Errorf("failed to send request: %v", err)
		return &RequestRet{
			respond: flow.Respond{
				Header: nil,
				Data:   []byte(fmt.Sprintf("failed to send request: %v", err)),
				Status: http.StatusServiceUnavailable,
			},
			selection: []int{0, 2},
		}, nil
	}

	header := make(map[string]interface{})
	for k, v := range response.Header {
		header[k] = v[0]
	}

	if response.StatusCode >= 100 && response.StatusCode < 400 {
		return &RequestRet{
			respond: flow.Respond{
				Header: header,
				Data:   response.Body,
				Status: response.StatusCode,
			},
			selection: []int{1, 2},
		}, nil
	}

	return &RequestRet{
		respond: flow.Respond{
			Header: header,
			Data:   response.Body,
			Status: response.StatusCode,
		},
		selection: []int{0, 2},
	}, nil
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

func (n *Request) IsRespond() bool {
	return false
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

func (n *Request) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Request) NextCount() int {
	return len(n.outputs)
}

func (n *Request) ActiveInput(node string) {
	for i := range n.inputs {
		if n.inputs[i].Node == node {
			if !n.inputs[i].Active {
				n.inputs[i].Active = true
				// input_1 for dynamic variable
				if n.inputs[i].InputName == flow.Input1 {
					n.lockCtx, n.lockCancel = context.WithCancel(context.Background())
				}
			}

			break
		}
	}
}

func (n *Request) CheckData() string {
	return ""
}

func (n *Request) Check() {
	n.checked = true
}

func (n *Request) IsChecked() bool {
	return n.checked
}

func NewRequest(_ context.Context, data flow.NodeData) flow.Noder {
	inputs := flow.PrepareInputs(data.Inputs)

	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

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
