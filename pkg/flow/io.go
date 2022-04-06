package flow

import "sort"

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

	for _, input := range inputs {
		for _, connection := range input.Connections {
			retInputs = append(retInputs, Inputs{Node: connection.Node})
		}
	}

	return retInputs
}
