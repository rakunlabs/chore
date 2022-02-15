package nodes

import (
	"context"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gorm.io/gorm"
)

var respondType = "respond"

// Respond node has one input.
type Respond struct {
	typeName string
}

// Run get values from active input nodes.
func (n *Respond) Run(_ context.Context, _ *registry.AppStore, value []byte) ([]byte, error) {
	return value, nil
}

func (n *Respond) GetType() string {
	return n.typeName
}

func (n *Respond) Fetch(ctx context.Context, db *gorm.DB) error {
	return nil
}

func (n *Respond) IsFetched() bool {
	return true
}

func (n *Respond) Validate() error {
	return nil
}

func (n *Respond) Next() []string {
	return nil
}

func (n *Respond) CheckData() string {
	return ""
}

func (n *Respond) ActiveInput(string) {}

func NewRespond(data flow.NodeData) flow.Noder {
	return &Respond{
		typeName: respondType,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[respondType] = NewRespond
}
