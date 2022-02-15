package models

import (
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"
)

type UserData struct {
	Name  string `json:"name" gorm:"unique;uniqueIndex;not null" example:"userX"`
	Email string `json:"email" example:"userx@worldline.com"`
	apimodels.Groups
}

type UserPrivate struct {
	Password string `json:"password" gorm:"not null" example:"pass1234"`
}

type UserPure struct {
	UserPrivate
	UserData
}

type User struct {
	UserPure
	apimodels.ModelCU
}
