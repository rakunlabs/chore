package models

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type TemplatePure struct {
	Name    string `json:"name" gorm:"uniqueIndex;not null" example:"deepcore/template1"`
	Content string `json:"content" swaggertype:"string" format:"base64" example:"aGVsbG8ge3submFtZX19Cg=="`
}

type Template struct {
	TemplatePure
	apimodels.ModelS
}

type FolderPure struct {
	Folder string `json:"folder" example:"deepcore/"`
	Item   string `json:"item" example:"template1" gorm:"not null"`
	Name   string `json:"name" example:"deepcore/template1" gorm:"uniqueIndex;not null"`
}

type Folder struct {
	FolderPure
	apimodels.ID
}
