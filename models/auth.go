package models

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"

	"gorm.io/datatypes"
)

type AuthPure struct {
	Name    string            `json:"name" gorm:"unique;uniqueIndex;not null" example:"jira-deepcore"`
	Headers datatypes.JSONMap `json:"headers" swaggertype:"object,string" example:"Content-Type:application/json"`
	apimodels.Groups
}

type Auth struct {
	AuthPure
	apimodels.ModelCU
}
