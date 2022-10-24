package js

import (
	"context"
	"fmt"
	"time"

	"github.com/dop251/goja"
	"gopkg.in/yaml.v3"
)

func toObject(v []byte) interface{} {
	var m interface{}

	if err := yaml.Unmarshal(v, &m); err != nil {
		panic(err)
	}

	return m
}

func toString(v []byte) string {
	return string(v)
}

func sleep(ctx context.Context) interface{} {
	return func(durationStr string) {
		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			panic(err)
		}

		timer := time.NewTimer(duration)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
		}
	}
}

func setValue(_ interface{}) {}

type commands struct {
	fn    interface{}
	fnCtx func(context.Context) interface{}
	name  string
}

var commandList = []commands{
	{
		fn:   toObject,
		name: "toObject",
	},
	{
		fn:   toString,
		name: "toString",
	},
	{
		fn:   setValue,
		name: "setValue",
	},
	{
		fnCtx: sleep,
		name:  "sleep",
	},
}

func setScriptFuncs(ctx context.Context, runner *goja.Runtime, additional map[string]interface{}) error {
	for _, command := range commandList {
		fn := command.fn
		if command.fnCtx != nil {
			fn = command.fnCtx(ctx)
		}

		// custom functions set
		if err := runner.Set(command.name, fn); err != nil {
			return fmt.Errorf("%s command cannot set: %w", command.name, err)
		}
	}

	for name, fn := range additional {
		if err := runner.Set(name, fn); err != nil {
			return fmt.Errorf("%s command cannot set: %w", name, err)
		}
	}

	return nil
}
