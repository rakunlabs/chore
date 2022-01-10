package kv

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/consul/api"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/store/inf"
)

type Consul struct {
	client       *api.Client
	kv           *api.KV
	writeOptions *api.WriteOptions
	queryOptions *api.QueryOptions
	basePath     string
}

func (c *Consul) Put(key string, value []byte) inf.ErrCrud {
	keyPath := path.Join(c.basePath, key)

	if keyPath[0] == '/' {
		keyPath = keyPath[1:]
	}

	p := &api.KVPair{
		Key:   keyPath,
		Value: value,
	}

	_, err := c.kv.Put(p, c.writeOptions)
	if err != nil {
		return inf.ErrOperation{Err: fmt.Errorf("consul put failed %w", err), Code: 500}
	}

	return nil
}

func (c *Consul) Post(key string, value []byte) inf.ErrCrud {
	if key == "" {
		return inf.ErrOperation{Err: fmt.Errorf("consul wrong key"), Code: 406}
	}

	searchPath := "/"
	keyArray := strings.Split(key, "/")

	if len(keyArray) != 1 {
		searchPath = strings.Join(keyArray[:len(keyArray)-1], "/")
	}

	items, err := c.List(searchPath)
	if err != nil {
		return c.Put(key, value)
	}

	checkPath := key
	if checkPath[0] == '/' {
		checkPath = checkPath[1:]
	}

	for _, item := range items {
		if item == checkPath {
			return inf.ErrOperation{Err: fmt.Errorf("consul key %s already exists", key), Code: 409}
		}
	}

	return c.Put(path.Join(c.basePath, key), value)
}

func (c *Consul) Get(key string) ([][]byte, inf.ErrCrud) {
	pairs, _, err := c.kv.List(path.Join(c.basePath, key), c.queryOptions)
	if err != nil {
		return nil, inf.ErrOperation{Err: fmt.Errorf("consul get failed %w", err), Code: 500}
	}

	if pairs == nil {
		return nil, inf.ErrOperation{Err: fmt.Errorf("consul cannot find %s", key), Code: 404}
	}

	collectValues := make([][]byte, 0, len(pairs))

	for i := range pairs {
		collectValues = append(collectValues, pairs[i].Value)
	}

	return collectValues, nil
}

func (c *Consul) List(key string) ([]string, inf.ErrCrud) {
	searchPath := path.Join(c.basePath, key) + "/"

	pairs, _, err := c.kv.Keys(searchPath, "/", c.queryOptions)
	if err != nil {
		return nil, inf.ErrOperation{Err: fmt.Errorf("consul get keys failed %w", err), Code: 500}
	}

	if pairs == nil {
		return nil, inf.ErrOperation{Err: fmt.Errorf("consul cannot find %s", searchPath), Code: 404}
	}

	for i := range pairs {
		pairs[i] = strings.TrimLeft(pairs[i], c.basePath)
	}

	return pairs, nil
}

func (c *Consul) Delete(key string) inf.ErrCrud {
	if key == "" {
		return inf.ErrOperation{Err: fmt.Errorf("consul wrong key"), Code: 406}
	}

	searchPath := "/"
	keyArray := strings.Split(key, "/")

	if len(keyArray) != 1 {
		searchPath = strings.Join(keyArray[:len(keyArray)-1], "/")
	}

	items, errCrud := c.List(searchPath)
	if errCrud != nil {
		return inf.ErrOperation{Err: fmt.Errorf("consul delete %q not found", key), Code: 404}
	}

	checkPath := key
	if checkPath[0] == '/' {
		checkPath = checkPath[1:]
	}

	found := false

	for _, item := range items {
		if item == checkPath {
			found = true
		}
	}

	if !found {
		return inf.ErrOperation{Err: fmt.Errorf("consul delete %q not found", key), Code: 404}
	}

	_, err := c.kv.Delete(path.Join(c.basePath, checkPath), c.writeOptions)
	if err != nil {
		return inf.ErrOperation{Err: fmt.Errorf("consul delete failed %w", err), Code: 500}
	}

	return nil
}

// func (c *Consul) Watch(ctx context.Context, wg *sync.WaitGroup, key string, fn func(interface{})) error {
// 	if fn == nil {
// 		return fmt.Errorf("nil function")
// 	}

// 	plan, err := watch.Parse(map[string]interface{}{
// 		"type": "key",
// 		"key":  key,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("consul watch failed %w", err)
// 	}

// 	plan.HybridHandler = func(_ watch.BlockingParamVal, raw interface{}) {
// 		if raw == nil {
// 			return
// 		}

// 		v, ok := raw.(*api.KVPair)
// 		if !ok {
// 			log.Ctx(ctx).Debug().Msg("incomprehensible value")

// 			return
// 		}

// 		fn(v.Value)
// 	}

// 	wg.Add(1)

// 	go func() {
// 		if err := plan.RunWithClientAndHclog(c.client, hclog.NewNullLogger()); err != nil {
// 			log.Error().Err(err).Msg("closed watch")
// 		}

// 		wg.Done()
// 	}()

// 	wg.Add(1)

// 	go func() {
// 		// wait context cancel
// 		<-ctx.Done()

// 		// close watch
// 		plan.Stop()

// 		wg.Done()
// 	}()

// 	return nil
// }

func (c *Consul) Close() error {
	return nil
}

func NewConsul(ctx context.Context, basePath string) (inf.CRUD, error) {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to consul %w", err)
	}

	consul := &Consul{
		client:       client,
		kv:           client.KV(),
		writeOptions: new(api.WriteOptions).WithContext(ctx),
		queryOptions: new(api.QueryOptions).WithContext(ctx),
		basePath:     strings.Trim(basePath, "/") + "/",
	}

	return consul, nil
}
