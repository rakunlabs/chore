package models

import (
	"github.com/rakunlabs/chore/pkg/models/apimodels"
)

type GroupPure struct {
	Name string `json:"name" gorm:"uniqueIndex;not null"`
}

type Group struct {
	GroupPure
	apimodels.ModelCU
}
