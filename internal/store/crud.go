package store

import (
	"context"
	"fmt"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/inf"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/kv"
)

func Define(ctx context.Context, store, basePath string) (inf.CRUD, error) {
	//nolint:gocritic
	switch store {
	case "consul":
		crud, err := kv.NewConsul(ctx, basePath)
		if err != nil {
			return nil, fmt.Errorf("newConsul %w", err)
		}

		return crud, nil
	}

	return nil, fmt.Errorf("not find store type %v", store)
}
