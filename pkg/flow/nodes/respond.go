package nodes

import (
	"context"
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/worldline-go/chore/pkg/flow"
	"github.com/worldline-go/chore/pkg/registry"
)

var respondType = "respond"

type RespondRet struct {
	respond flow.Respond
}

func (r *RespondRet) GetBinaryData() []byte {
	return r.respond.Data
}

func (r *RespondRet) GetRespond() flow.Respond {
	return r.respond
}

var _ flow.NodeRetRespond = &RespondRet{}

// Respond node has one input.
type Respond struct {
	statusCodeRaw string
	headersRaw    string
	getData       bool
	checked       bool
}

// Run get values from active input nodes.
func (n *Respond) Run(ctx context.Context, _ *registry.AppStore, value flow.NodeRet, _ string) (flow.NodeRet, error) {
	var headers map[string]interface{}
	if err := yaml.Unmarshal([]byte(n.headersRaw), &headers); err != nil {
		return nil, fmt.Errorf("faild unmarshal headers in request: %v", err)
	}

	if n.getData {
		if vRespond, ok := value.(flow.NodeRetRespondData); ok {
			// combine headers
			respond := vRespond.GetRespondData()
			if respond.Header == nil {
				respond.Header = headers
			} else {
				for k, v := range headers {
					respond.Header[k] = v
				}
			}

			return &RespondRet{
				respond: vRespond.GetRespondData(),
			}, nil
		}
	}

	statusCode, err := strconv.Atoi(n.statusCodeRaw)
	if err != nil {
		log.Ctx(ctx).Warn().Msgf("status code %v cannot convert to integer, passing with 200", n.statusCodeRaw)

		statusCode = 200
	}

	return &RespondRet{
		respond: flow.Respond{
			Header: headers,
			Data:   value.GetBinaryData(),
			Status: statusCode,
		},
	}, nil
}

func (n *Respond) GetType() string {
	return respondType
}

func (n *Respond) Fetch(ctx context.Context, db *gorm.DB) error {
	return nil
}

func (n *Respond) IsFetched() bool {
	return true
}

func (n *Respond) IsRespond() bool {
	return true
}

func (n *Respond) Validate() error {
	return nil
}

func (n *Respond) Next(int) []flow.Connection {
	return nil
}

func (n *Respond) NextCount() int {
	return 0
}

func (n *Respond) Check() {
	n.checked = true
}

func (n *Respond) IsChecked() bool {
	return n.checked
}

func (n *Respond) ActiveInput(string) {}

func NewRespond(_ context.Context, data flow.NodeData) flow.Noder {
	headersRaw, _ := data.Data["headers"].(string)

	statusCodeRaw, _ := data.Data["status"].(string)

	// getData "true" or "false"
	getData, _ := data.Data["get"].(string)

	return &Respond{
		statusCodeRaw: statusCodeRaw,
		headersRaw:    headersRaw,
		getData:       getData == "true",
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[respondType] = NewRespond
}
