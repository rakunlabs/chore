package models

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type EmailPure struct {
	Host  string `json:"host"`
	Email string `json:"email"`
	Port  int    `json:"port"`
}

type EmailPrivate struct {
	Password string `json:"password"`
}

type Email struct {
	EmailPrivate
	EmailPure
}

type SettingsPure struct {
	Namespace string `json:"namespace" gorm:"unique;uniqueIndex;not null" example:"application"`
	Email     `json:"email"`
}

type Settings struct {
	apimodels.ModelCUPure
	SettingsPure
}
