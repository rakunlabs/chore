package apimodels

import "errors"

var (
	ErrRequiredID     = errors.New("id is required")
	ErrRequiredName   = errors.New("name is required")
	ErrRequiredIDName = errors.New("required at leats one of id or name")
	ErrNotFound       = errors.New("not found any releated data")
)
