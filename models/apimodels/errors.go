package apimodels

import "errors"

var (
	ErrRequiredID     = errors.New("id is required")
	ErrRequiredIDName = errors.New("required at leats one of id or name")
)
