package nodes

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/email"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
	"gopkg.in/yaml.v3"

	"gorm.io/gorm"
)

var emailType = "email"

// Email node has one input.
type Email struct {
	client   email.Client
	typeName string
	fetched  bool
	headers  map[string][]string
}

// Run get values from active input nodes.
func (n *Email) Run(_ context.Context, _ *registry.AppStore, value []byte, _ string) ([]byte, error) {
	if err := n.client.Send(value, n.headers); err != nil {
		return nil, fmt.Errorf("email send failed: %v", err)
	}

	return value, nil
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

func (n *Email) Next() []flow.Connection {
	return nil
}

func (n *Email) CheckData() string {
	return ""
}

func (n *Email) ActiveInput(string) {}

func NewEmail(data flow.NodeData) flow.Noder {
	var headers map[string][]string
	if err := yaml.Unmarshal([]byte(data.Data["script"].(string)), &headers); err != nil {
		log.Error().Err(err).Msg("email yaml unmarshal failed")
	}

	return &Email{
		typeName: emailType,
		headers:  headers,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[emailType] = NewEmail
}
