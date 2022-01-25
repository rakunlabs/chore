package models

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"

	"gorm.io/datatypes"
)

type AuthPure struct {
	Name    string            `json:"name" gorm:"uniqueIndex;not null" example:"jira-deepcore"`
	Headers datatypes.JSONMap `json:"header" swaggertype:"object,string" example:"Content-Type:application/json"`
	URL     string            `json:"url" example:"http://localhost:9090"`
	Method  string            `json:"method" example:"POST"`
}

type Auth struct {
	AuthPure
	apimodels.ModelS
}
