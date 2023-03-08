package nodes

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/rytsh/liz/utils/templatex"
	"github.com/worldline-go/chore/models"
	"github.com/worldline-go/chore/pkg/flow"
	"github.com/worldline-go/chore/pkg/flow/convert"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/request"
	"github.com/worldline-go/chore/pkg/transfer"

	"github.com/rs/zerolog"
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

type retryRaw struct {
	Codes   string
	DeCodes string
}

type renderedValues struct {
	url           string
	addHeadersRaw string
	method        string
}

// Request node has one input and one output.
type Request struct {
	reg           *flow.NodesReg
	nodeID        string
	lockCtx       context.Context
	lockCancel    context.CancelFunc
	headers       map[string]interface{}
	retry         *request.Retry
	retryRaw      retryRaw
	url           string
	addHeadersRaw string
	method        string
	auth          string
	outputs       [][]flow.Connection
	inputs        []flow.Inputs
	inputHolder   inputHolderRequest
	mutex         sync.Mutex
	fetched       bool
	checked       bool
	disabled      bool
	payloadNil    bool
	skipVerify    bool
	poolClient    bool
	stuckContext  context.Context
	log           *zerolog.Logger
	client        *request.Client
	tags          []string
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Request) Run(ctx context.Context, _ *sync.WaitGroup, reg *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
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

	// check it has value
	var useValues []byte
	if vRet, ok := value.(flow.NodeRetValues); ok {
		useValues = vRet.GetBinaryValues()
	}

	if useValues == nil && n.lockCtx != nil {
		// increase count
		n.reg.UpdateStuck(flow.CountStuckIncrease, false)

		select {
		case <-n.lockCtx.Done():
			// continue process
		default:
			// these events not happen at same time mostly
			select {
			case <-n.stuckContext.Done():
				return nil, fmt.Errorf("stuck detected, terminated node wait")
			case <-ctx.Done():
				log.Ctx(ctx).Warn().Msg("program closed, terminated node wait")

				return nil, flow.ErrStopGoroutine
			case <-n.lockCtx.Done():
				// continue process
			}
		}

		n.reg.UpdateStuck(flow.CountStuckDecrease, true)
	}

	// check value and render it
	rendered := renderedValues{
		url:           n.url,
		addHeadersRaw: n.addHeadersRaw,
		method:        n.method,
	}

	var requestValues interface{}
	if useValues != nil {
		requestValues = transfer.BytesToData(useValues)
	} else {
		requestValues = transfer.BytesToData(n.inputHolder.value)
	}

	if requestValues != nil {
		// render url
		renderedValue, err := reg.Template.ExecuteBuffer(templatex.WithData(requestValues), templatex.WithContent(n.url))
		if err != nil {
			return nil, fmt.Errorf("template url cannot render: %v", err)
		}

		rendered.url = string(renderedValue)

		// render method
		renderedValue, err = reg.Template.ExecuteBuffer(templatex.WithData(requestValues), templatex.WithContent(n.method))
		if err != nil {
			return nil, fmt.Errorf("template method cannot render: %v", err)
		}

		rendered.method = string(renderedValue)

		// render headers
		renderedValue, err = reg.Template.ExecuteBuffer(templatex.WithData(requestValues), templatex.WithContent(n.addHeadersRaw))
		if err != nil {
			return nil, fmt.Errorf("template headers cannot render: %v", err)
		}

		rendered.addHeadersRaw = string(renderedValue)
	}

	var addHeaders map[string]interface{}
	if err := yaml.Unmarshal([]byte(rendered.addHeadersRaw), &addHeaders); err != nil {
		return nil, fmt.Errorf("faild unmarshal headers in request: %v", err)
	}

	headers := make(map[string]interface{}, len(n.headers)+len(addHeaders)+1)
	if v, _ := ctx.Value("request_id").(string); v != "" {
		headers["X-Request-Id"] = v
	}
	for k := range n.headers {
		headers[k] = n.headers[k]
	}

	for k := range addHeaders {
		headers[k] = addHeaders[k]
	}

	var payload []byte
	if !n.payloadNil {
		payload = value.GetBinaryData()
	}

	if n.client == nil {
		return nil, fmt.Errorf("http client not set")
	}

	response, err := n.client.Send(
		ctx,
		rendered.url,
		rendered.method,
		headers,
		payload,
		n.retry,
		n.skipVerify,
	)
	if err != nil {
		// return nil, fmt.Errorf("failed to send request: %v", err)
		return &RequestRet{
			respond: flow.Respond{
				Header: nil,
				Data:   []byte(fmt.Sprint(err)),
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
	return requestType
}

func (n *Request) Fetch(ctx context.Context, db *gorm.DB) error {
	if n.auth == "" {
		n.fetched = true

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

func (n *Request) Validate(ctx context.Context) error {
	// correction
	if n.method == "" {
		n.method = "POST"
	}

	if n.url == "" {
		return fmt.Errorf("url is empty")
	}

	// fill retry values
	retryCodes, err := getCodes(n.retryRaw.Codes)
	if err != nil {
		return err
	}

	retryDeCodes, err := getCodes(n.retryRaw.DeCodes)
	if err != nil {
		return err
	}

	n.retry = &request.Retry{
		EnabledStatusCodes:  retryCodes,
		DisabledStatusCodes: retryDeCodes,
	}

	n.client = request.NewClient(request.Config{
		SkipVerify: n.skipVerify,
		Pooled:     n.poolClient,
		Log:        n.log,
	})

	n.stuckContext = n.reg.GetStuctCancel(ctx)

	return nil
}

func (n *Request) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Request) NextCount() int {
	return len(n.outputs)
}

func (n *Request) IsDisabled() bool {
	return n.disabled
}

func (n *Request) ActiveInput(node string, tags map[string]struct{}) {
	if !convert.IsTagsEnabled(n.tags, tags) {
		n.disabled = true

		return
	}

	for i := range n.inputs {
		if n.inputs[i].Node == node {
			if !n.inputs[i].Active {
				n.inputs[i].Active = true
				// input_1 for dynamic variable
				if n.inputs[i].InputName == flow.Input1 {
					n.lockCtx, n.lockCancel = context.WithCancel(context.Background())
					n.reg.AddCleanup(n.lockCancel)
				}
			}
		}
	}
}

func (n *Request) Check() {
	n.checked = true
}

func (n *Request) IsChecked() bool {
	return n.checked
}

func (n *Request) NodeID() string {
	return n.nodeID
}

func (n *Request) Tags() []string {
	return n.tags
}

func NewRequest(ctx context.Context, reg *flow.NodesReg, data flow.NodeData, nodeID string) (flow.Noder, error) {
	inputs := flow.PrepareInputs(data.Inputs)

	// add outputs with order
	outputs := flow.PrepareOutputs(data.Outputs)

	auth, _ := data.Data["auth"].(string)
	method, _ := data.Data["method"].(string)
	url, _ := data.Data["url"].(string)
	addHeadersRaw, _ := data.Data["headers"].(string)

	retryCodes, _ := data.Data["retry_codes"].(string)
	retryDeCodes, _ := data.Data["retry_decodes"].(string)

	skipVerify := convert.GetBoolean(data.Data["skip_verify"])
	payloadNil := convert.GetBoolean(data.Data["payload_nil"])
	poolClient := convert.GetBoolean(data.Data["pool_client"])

	tags := convert.GetList(data.Data["tags"])

	l := log.Ctx(ctx).With().Str("component", requestType).Logger()

	return &Request{
		reg:           reg,
		inputs:        inputs,
		outputs:       outputs,
		auth:          auth,
		method:        method,
		url:           url,
		addHeadersRaw: addHeadersRaw,
		retryRaw: retryRaw{
			Codes:   strings.ReplaceAll(retryCodes, ",", " "),
			DeCodes: strings.ReplaceAll(retryDeCodes, ",", " "),
		},
		skipVerify: skipVerify,
		payloadNil: payloadNil,
		poolClient: poolClient,
		log:        &l,
		nodeID:     nodeID,
		tags:       tags,
	}, nil
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[requestType] = NewRequest
}

func getCodes(codes string) ([]int, error) {
	if codes == "" {
		return nil, nil
	}

	var retryCodes []int

	for _, code := range strings.Fields(codes) {
		codeInt, err := strconv.Atoi(code)
		if err != nil {
			return nil, fmt.Errorf("value %s cannot convert to integer", code)
		}

		retryCodes = append(retryCodes, codeInt)
	}

	return retryCodes, nil
}
