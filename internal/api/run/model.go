package run

import "time"

var defaultTimeout = 30 * time.Second

type runModel struct {
	Script   []byte   `json:"script"`
	Settings settings `json:"settings"`
	Inputs   []byte   `json:"inputs"`
}

func defaultRunModel() runModel {
	return runModel{
		Settings: settings{
			TimeoutDuration: defaultTimeout,
		},
	}
}

type settings struct {
	Timeout         string        `json:"timeout"`
	TimeoutDuration time.Duration `json:"-"`
	Async           bool          `json:"async"`
}
