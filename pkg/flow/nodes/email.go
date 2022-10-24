package nodes

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/worldline-go/chore/models"
	"github.com/worldline-go/chore/pkg/email"
	"github.com/worldline-go/chore/pkg/flow"
	"github.com/worldline-go/chore/pkg/registry"
	"github.com/worldline-go/chore/pkg/transfer"
)

var emailType = "email"

type inputHolderEmail struct {
	value []byte
	exist bool
}

type EmailRet struct {
	output []byte
}

func (r *EmailRet) GetBinaryData() []byte {
	return r.output
}

// Email node has one input.
type Email struct {
	reg         *flow.NodesReg
	lockCtx     context.Context
	lockCancel  context.CancelFunc
	values      map[string]string
	client      email.Client
	inputs      []flow.Inputs
	inputHolder inputHolderEmail
	mutex       sync.Mutex
	fetched     bool
	checked     bool
	nodeID      string
}

// Run get values from active input nodes.
func (n *Email) Run(ctx context.Context, _ *sync.WaitGroup, reg *registry.AppStore, value flow.NodeRet, input string) (flow.NodeRet, error) {
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
		case <-n.reg.GetStuckCtx().Done():
			return nil, fmt.Errorf("stuck detected, terminated node request")
		case <-ctx.Done():
			log.Ctx(ctx).Warn().Msg("program closed, terminated node request")

			return nil, flow.ErrStopGoroutine
		case <-n.lockCtx.Done():
			// continue process
		}

		n.reg.UpdateStuck(flow.CountStuckDecrease, true)
	}

	// check value and render it
	headers := make(map[string][]string)

	var requestValues interface{}
	if useValues != nil {
		requestValues = transfer.BytesToData(useValues)
	} else {
		requestValues = transfer.BytesToData(n.inputHolder.value)
	}

	for key, value := range n.values {
		payload := value

		if requestValues != nil {
			// render
			rendered, err := reg.Template.Execute(requestValues, value)
			if err != nil {
				return nil, fmt.Errorf("template cannot render: %v", err)
			}

			payload = rendered
		}

		if key == "Subject" {
			headers[key] = []string{payload}

			continue
		}

		payloadSlice := strings.Fields(strings.ReplaceAll(payload, ",", " "))
		if len(payloadSlice) > 0 {
			headers[key] = payloadSlice
		}
	}

	if err := n.client.Send(value.GetBinaryData(), headers); err != nil {
		return nil, fmt.Errorf("failed to send email: values %v, err %v", headers, err)
	}

	return &EmailRet{output: value.GetBinaryData()}, nil
}

func (n *Email) GetType() string {
	return emailType
}

func (n *Email) Fetch(ctx context.Context, db *gorm.DB) error {
	getData := models.Email{}

	query := db.WithContext(ctx).Model(&models.Settings{}).Where("namespace = ?", "application")
	result := query.First(&getData)

	if result.Error != nil {
		return fmt.Errorf("email fetch failed: %v", result.Error)
	}

	n.client = email.NewClient(getData.Host, getData.Port, getData.NoAuth, getData.Email, getData.Password)

	if n.values["From"] == "" {
		n.values["From"] = getData.Email
	}

	n.fetched = true

	return nil
}

func (n *Email) IsFetched() bool {
	return n.fetched
}

func (n *Email) IsRespond() bool {
	return false
}

func (n *Email) Validate() error {
	return nil
}

func (n *Email) Next(int) []flow.Connection {
	return nil
}

func (n *Email) NextCount() int {
	return 0
}

func (n *Email) Check() {
	n.checked = true
}

func (n *Email) IsChecked() bool {
	return n.checked
}

func (n *Email) ActiveInput(node string) {
	for i := range n.inputs {
		if n.inputs[i].Node == node {
			if !n.inputs[i].Active {
				n.inputs[i].Active = true
				// input_1 for dynamic variable
				if n.inputs[i].InputName == flow.Input1 {
					n.lockCtx, n.lockCancel = context.WithCancel(context.Background())
				}
			}
		}
	}
}

func (n *Email) NodeID() string {
	return n.nodeID
}

func NewEmail(_ context.Context, reg *flow.NodesReg, data flow.NodeData, nodeID string) (flow.Noder, error) {
	inputs := flow.PrepareInputs(data.Inputs)

	values := make(map[string]string)

	values["From"], _ = data.Data["from"].(string)
	values["To"], _ = data.Data["to"].(string)
	values["Cc"], _ = data.Data["cc"].(string)
	values["Bcc"], _ = data.Data["bcc"].(string)
	values["Subject"], _ = data.Data["subject"].(string)

	return &Email{
		reg:    reg,
		values: values,
		inputs: inputs,
		nodeID: nodeID,
	}, nil
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[emailType] = NewEmail
}
