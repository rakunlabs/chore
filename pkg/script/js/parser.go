package js

import "gopkg.in/yaml.v3"

func ParseInputs(data []byte) []interface{} {
	var inputs []interface{}

	if err := yaml.Unmarshal(data, &inputs); err != nil {
		return []interface{}{string(data)}
	}

	return inputs
}
