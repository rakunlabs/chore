package nodes

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/email"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

var emailType = "email"

type inputHolderEmail struct {
	value []byte
	data  []byte
}

// Email node has one input.
type Email struct {
	values      map[string]string
	client      email.Client
	typeName    string
	inputHolder inputHolderEmail
	wait        int
	lock        sync.Mutex
	fetched     bool
}

// Run get values from active input nodes.
func (n *Email) Run(_ context.Context, reg *registry.AppStore, value []byte, input string) ([][]byte, error) {
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

	headers := make(map[string][]string)

	// check value and render it
	if n.inputHolder.value != nil {
		var requestValues map[string]interface{}
		if err := yaml.Unmarshal(n.inputHolder.value, &requestValues); err != nil {
			return nil, fmt.Errorf("failed to unmarshal: %v", err)
		}

		for key, value := range n.values {
			// render
			payload, err := reg.Template.Ext(requestValues, value)
			if err != nil {
				return nil, fmt.Errorf("template cannot render: %v", err)
			}

			headers[key] = strings.Split(string(payload), ",")
		}
	}

	if err := n.client.Send(value, headers); err != nil {
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	return [][]byte{value}, nil
}

func (n *Email) GetType() string {
	return n.typeName
}

func (n *Email) Fetch(ctx context.Context, db *gorm.DB) error {
	getData := models.Email{}

	query := db.WithContext(ctx).Model(&models.Settings{}).Where("namespace = ?", "application")
	result := query.First(&getData)

	if result.Error != nil {
		return fmt.Errorf("email fetch failed: %v", result.Error)
	}

	n.client = email.NewClient(getData.Host, getData.Port, getData.Email, getData.Password)

	n.fetched = true

	return nil
}

func (n *Email) IsFetched() bool {
	return n.fetched
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

func (n *Email) CheckData() string {
	return ""
}

func (n *Email) ActiveInput(string) {}

func NewEmail(data flow.NodeData) flow.Noder {
	values := make(map[string]string)

	values["From"], _ = data.Data["from"].(string)
	values["To"], _ = data.Data["to"].(string)
	values["Cc"], _ = data.Data["cc"].(string)
	values["Bcc"], _ = data.Data["bcc"].(string)
	values["Subject"], _ = data.Data["subject"].(string)

	return &Email{
		typeName: emailType,
		values:   values,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[emailType] = NewEmail
}
