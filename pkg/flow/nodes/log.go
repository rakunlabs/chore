package nodes

import (
	"context"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/flow"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var logType = "log"

// Respond node has one input.
type Log struct {
	typeName  string
	message   string
	outputs   [][]flow.Connection
	printData bool
	logLevel  zerolog.Level
}

// Run get values from everywhere no need to check active input.
func (n *Log) Run(ctx context.Context, _ *registry.AppStore, value []byte, _ string) ([][]byte, error) {
	var logEvent *zerolog.Event

	switch n.logLevel {
	case zerolog.DebugLevel:
		logEvent = log.Ctx(ctx).Debug()
	case zerolog.InfoLevel:
		logEvent = log.Ctx(ctx).Info()
	case zerolog.WarnLevel:
		logEvent = log.Ctx(ctx).Warn()
	case zerolog.ErrorLevel:
		logEvent = log.Ctx(ctx).Error()
	default:
		logEvent = log.Ctx(ctx).Debug()
	}

	if n.printData {
		logEvent = logEvent.Str("data", string(value))
	}

	logEvent.Msg(n.message)

	return [][]byte{value}, nil
}

func (n *Log) GetType() string {
	return n.typeName
}

func (n *Log) Fetch(ctx context.Context, db *gorm.DB) error {
	return nil
}

func (n *Log) IsFetched() bool {
	return true
}

func (n *Log) Validate() error {
	return nil
}

func (n *Log) Next(i int) []flow.Connection {
	return n.outputs[i]
}

func (n *Log) NextCount() int {
	return len(n.outputs)
}

func (n *Log) CheckData() string {
	return ""
}

func (n *Log) ActiveInput(string) {}

func NewLog(data flow.NodeData) flow.Noder {
	outputs := flow.PrepareOutputs(data.Outputs)

	// printData "true" or "false"
	printData, _ := data.Data["data"].(string)
	// loglevel "debug", "info", "warn", "error"
	level, _ := data.Data["level"].(string)
	// message default ""
	message, _ := data.Data["message"].(string)

	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.DebugLevel
	}

	return &Log{
		typeName:  logType,
		outputs:   outputs,
		printData: printData == "true",
		logLevel:  logLevel,
		message:   message,
	}
}

//nolint:gochecknoinits // moduler nodes
func init() {
	flow.NodeTypes[logType] = NewLog
}
