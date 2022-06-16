package nodes

import (
	"context"
	"strconv"
	"strings"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gorm.io/gorm"
)

var endpointType = "endpoint"

type EndpointRet struct {
	output []byte
}

func (r *EndpointRet) GetBinaryData() []byte {
	return r.output
}

// Endpoint node has one output.
type Endpoint struct {
	endpoint string
	outputs  [][]flow.Connection
	methods  []string
	checked  bool
	public   bool
}

var _ flow.NoderEndpoint = &Endpoint{}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Endpoint) Run(_ context.Context, _ *registry.AppStore, value flow.NodeRet, _ string) (flow.NodeRet, error) {
	return &EndpointRet{output: value.GetBinaryData()}, nil
}

func (n *Endpoint) GetType() string {
	return endpointType
}

func (n *Endpoint) Fetch(ctx context.Context, db *gorm.DB) error {
	return nil
}

func (n *Endpoint) IsFetched() bool {
	return true
}

func (n *Endpoint) IsRespond() bool {
	return false
}

func (n *Endpoint) Validate() error {
	return nil
}

func (n *Endpoint) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Endpoint) NextCount() int {
	return len(n.outputs)
}

func (n *Endpoint) Check() {
	n.checked = true
}

func (n *Endpoint) IsChecked() bool {
	return n.checked
}

func (n *Endpoint) ActiveInput(string) {}

func (n *Endpoint) Endpoint() string {
	return n.endpoint
}

func (n *Endpoint) Methods() []string {
	return n.methods
}

func NewEndpoint(_ context.Context, data flow.NodeData) flow.Noder {
	outputs := flow.PrepareOutputs(data.Outputs)

	endpoint, _ := data.Data["endpoint"].(string)
	methodsRaw, _ := data.Data["methods"].(string)
	publicRaw, _ := data.Data["public"].(string)

	methodsRaw = strings.ReplaceAll(methodsRaw, " ", "")
	publicRaw = strings.ReplaceAll(publicRaw, " ", "")

	methods := strings.Split(methodsRaw, ",")
	public, _ := strconv.ParseBool(publicRaw)

	return &Endpoint{
		outputs:  outputs,
		endpoint: endpoint,
		methods:  methods,
		public:   public,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[endpointType] = NewEndpoint
}
