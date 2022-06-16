package flow

import (
	"context"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/registry"
)

func StartFlow(
	ctx context.Context,
	controlName, endPoint, method string,
	content []byte,
	appStore *registry.AppStore,
	value []byte,
) (*NodesReg, error) {
	nodesData, err := ParseData(content)
	if err != nil {
		return nil, err
	}

	nodesReg, err := DataToNode(ctx, controlName, endPoint, method, nodesData, appStore)
	if err != nil {
		return nil, err
	}

	if err := VisitAndFetch(nodesReg); err != nil {
		return nil, err
	}

	nodesReg.wg.Add(1)

	go GoAndRun(nodesReg, value)

	return nodesReg, nil
}
