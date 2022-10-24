package transfer

import (
	"gopkg.in/yaml.v3"
)

func BytesToData(data []byte) interface{} {
	// check if data is nil
	if data == nil {
		return nil
	}

	var vX interface{}

	if err := yaml.Unmarshal(data, &vX); err == nil {
		return vX
	}

	return string(data)
}
