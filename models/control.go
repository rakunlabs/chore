package models

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type ControlPure struct {
	Name string `json:"name" gorm:"unique;uniqueIndex;not null"`
	apimodels.Groups
}

type ControlPureContent struct {
	Content string `json:"content" swaggertype:"string" format:"base64" example:"aGVsbG8ge3submFtZX19Cg=="`
	ControlPure
}

type Control struct {
	ControlPureContent
	apimodels.ModelCU
}
