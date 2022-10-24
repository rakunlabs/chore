package registry

import "sync"

type Option func(*Registry)

func WithWaitGroup(wg *sync.WaitGroup) Option {
	return func(r *Registry) {
		r.WG = wg
	}
}
