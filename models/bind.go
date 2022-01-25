package models

import (
	"github.com/google/uuid"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type BindPure struct {
	AuthID     *uuid.UUID `json:"auth_id"`
	TemplateID *uuid.UUID `json:"template_id"`
	Name       string     `json:"name" gorm:"uniqueIndex;not null"`
}

type Bind struct {
	BindPure
	// AuthID     *uuid.UUID `json:"auth_id"`
	// TemplateID *uuid.UUID `json:"template_id"`
	Template Template `json:"template" gorm:"foreignKey:TemplateID"`
	apimodels.ModelS
	// Name string `json:"name" gorm:"uniqueIndex;not null"`
	Auth Auth `json:"auth" gorm:"foreignKey:AuthID"`
}
