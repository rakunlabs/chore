package models

import "gorm.io/gorm"

type TestPure struct {
	Name string `json:"name" gorm:"uniqueIndex;not null" example:"userX"`
}

type Test struct {
	gorm.Model
	TestPure
}
