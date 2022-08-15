package models

import (
	"github.com/worldline-go/chore/models/apimodels"

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
