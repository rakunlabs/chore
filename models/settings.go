package models

import (
	"github.com/worldline-go/chore/models/apimodels"
)

type EmailPure struct {
	Host   string `json:"host"`
	Email  string `json:"email"`
	Port   int    `json:"port"`
	NoAuth bool   `json:"no_auth"`
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
