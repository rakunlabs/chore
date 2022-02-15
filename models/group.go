package models

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type GroupPure struct {
	Name string `json:"name" gorm:"uniqueIndex;not null"`
}

type Group struct {
	GroupPure
	apimodels.ModelCU
}
