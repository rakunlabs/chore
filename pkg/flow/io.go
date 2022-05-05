package flow

import (
	"sort"
)

var Input1 = "input_1"

func PrepareOutputs(outputs NodeConnection) [][]Connection {
	orderKey := make([]string, 0, len(outputs))
	for key := range outputs {
		orderKey = append(orderKey, key)
	}

	sort.Strings(orderKey)

	// add outputs with order
	retOutputs := make([][]Connection, 0, len(outputs))
	for _, key := range orderKey {
		retOutputs = append(retOutputs, outputs[key].Connections)
	}

	return retOutputs
}

func PrepareInputs(inputs NodeConnection) []Inputs {
	retInputs := make([]Inputs, 0, len(inputs))

	for inputName, input := range inputs {
		for _, connection := range input.Connections {
			retInputs = append(retInputs, Inputs{Node: connection.Node, InputName: inputName})
		}
	}

	return retInputs
}
