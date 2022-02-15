package apimodels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Limit for default limit api parameter.
var Limit = 20

type Meta struct {
	Limit  int   `json:"limit" query:"limit" example:"20"`
	Offset int   `json:"offset" query:"offset" example:"0"`
	Count  int64 `json:"count" query:"count" example:"35"`
}

type Error struct {
	Error interface{} `json:"error,omitempty" example:"some problem" swaggertype:"string"`
}

type Data struct {
	Data interface{} `json:"data,omitempty" swaggertype:"object,string"`
}

type DataMeta struct {
	Data
	Meta interface{} `json:"meta,omitempty"`
}

type API struct {
	DataMeta
	Error
}

// type Model struct {
// 	CreatedAt time.Time      `json:"created_at"`
// 	UpdatedAt time.Time      `json:"updated_at"`
// 	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
// 	ID
// }

type ModelCU struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ID
}

type ModelC struct {
	CreatedAt time.Time `json:"created_at"`
	ID
}

type ID struct {
	ID uuid.UUID `json:"id" gorm:"primarykey,type:uuid" example:"cf8a07d4-077e-402e-a46b-ac0ed50989ec"`
}

type Groups struct {
	Groups datatypes.JSON `json:"groups" example:"admin,group1" swaggertype:"array,string"`
}
