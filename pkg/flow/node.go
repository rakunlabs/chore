package flow

import (
	"context"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"gorm.io/gorm"
)

// NodeTypes hold new function of releated node.
var NodeTypes = make(map[string]func(NodeData) Noder)

type NodeData struct {
	Data    map[string]interface{} `json:"data"`
	Inputs  NodeConnection         `json:"inputs"`
	Outputs NodeConnection         `json:"outputs"`
	Name    string
}

type NodeConnection = map[string]Connections

type Connections = struct {
	Connections []Connection `json:"connections"`
}

type Connection = struct {
	Node   string `json:"node"`
	Output string `json:"output"`
}

// NodesData is content's represent.
type NodesData = map[string]NodeData

type Inputs struct {
	Node   string
	Active bool
}

// Noder for nodes like script, endpoint.
type Noder interface {
	GetType() string
	Run(context.Context, *registry.AppStore, []byte, string) ([]byte, error)
	Fetch(context.Context, *gorm.DB) error
	IsFetched() bool
	Validate() error
	ActiveInput(string)
	Next() []Connection
	CheckData() string
}
