package models

import (
	"github.com/worldline-go/chore/models/apimodels"
)

type GroupPure struct {
	Name string `json:"name" gorm:"uniqueIndex;not null"`
}

type Group struct {
	GroupPure
	apimodels.ModelCU
}
