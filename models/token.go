package models

import (
	"time"
)

var (
	TypePersonalAccessToken = "PAT"
	TypeAccessToken         = "X"
)

type TokenPure struct {
	Token string `json:"token" gorm:"primarykey"`
}

type Token struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	TokenPure
}
