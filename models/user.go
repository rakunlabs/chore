package models

import "gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models/apimodels"

type UserPure struct {
	Name     string `json:"name" gorm:"uniqueIndex;not null" example:"userX"`
	Password string `json:"password" example:"pass1234"`
	Admin    bool   `json:"admin" example:"true"`
}

type User struct {
	UserPure
	apimodels.ModelS
}
