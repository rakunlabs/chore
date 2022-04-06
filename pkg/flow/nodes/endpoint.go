package nodes

import (
	"context"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gorm.io/gorm"
)

var endpointType = "endpoint"

// Endpoint node has one output.
type Endpoint struct {
	endpoint string
	typeName string
	outputs  [][]flow.Connection
}

// Run get values from active input nodes and it will not run until last input comes.
func (n *Endpoint) Run(_ context.Context, _ *registry.AppStore, value []byte, _ string) ([][]byte, error) {
	return [][]byte{value}, nil
}

func (n *Endpoint) GetType() string {
	return n.typeName
}

func (n *Endpoint) Fetch(ctx context.Context, db *gorm.DB) error {
	return nil
}

func (n *Endpoint) IsFetched() bool {
	return true
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

func (n *Endpoint) CheckData() string {
	return n.endpoint
}

func (n *Endpoint) ActiveInput(string) {}

func NewEndpoint(data flow.NodeData) flow.Noder {
	outputs := flow.PrepareOutputs(data.Outputs)

	endpoint, _ := data.Data["endpoint"].(string)

	return &Endpoint{
		typeName: endpointType,
		outputs:  outputs,
		endpoint: endpoint,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[endpointType] = NewEndpoint
}
