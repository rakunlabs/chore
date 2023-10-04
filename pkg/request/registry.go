package request

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/worldline-go/auth"
	"github.com/worldline-go/auth/providers"
)

var ErrClientIDNotFound = errors.New("client id not found")

var (
	DefaultTimeDuration   = 10 * time.Minute
	DefaultTickerDuration = time.Minute
)

type Registry struct {
	clientIDs map[string]Auth
	timers    map[string]Timer
	ctx       context.Context //nolint:containedctx // application context

	mutex sync.RWMutex
}

var GlobalRegistry *Registry

type Auth struct {
	RoundTripper func(_ context.Context, transport http.RoundTripper) (http.RoundTripper, error)
	Cancel       context.CancelFunc
}

type Timer struct {
	Timer    *time.Timer
	ClientID string
}

func InitGlobalRegistry(ctx context.Context) *Registry {
	reg := &Registry{}

	reg.clientIDs = map[string]Auth{}
	reg.timers = map[string]Timer{}
	reg.ctx = ctx

	GlobalRegistry = reg

	return GlobalRegistry
}

func (r *Registry) AddService(cfg AuthConfig) (func(_ context.Context, transport http.RoundTripper) (http.RoundTripper, error), error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if cfg.ClientID == "" {
		return nil, ErrClientIDNotFound
	}

	if _, ok := r.clientIDs[cfg.ClientID]; ok {
		r.resetTimer(cfg.ClientID)
		// log.Debug().Msgf("using client id: %s", cfg.ClientID)

		return r.clientIDs[cfg.ClientID].RoundTripper, nil
	}

	log.Debug().Msgf("adding client id: %s", cfg.ClientID)

	provider := auth.Provider{Generic: &providers.Generic{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenURL,
		Scopes:       cfg.Scopes,
	}}
	providerActive := provider.ActiveProvider()

	ctx, cancel := context.WithCancel(r.ctx)
	shared, err := providerActive.NewOauth2Shared(ctx)
	if err != nil {
		cancel()

		return nil, err //nolint:wrapcheck // no need
	}

	timer := time.NewTimer(DefaultTimeDuration)

	r.clientIDs[cfg.ClientID] = Auth{
		RoundTripper: shared.RoundTripper,
		Cancel:       cancel,
	}

	r.timers[cfg.ClientID] = Timer{
		Timer:    timer,
		ClientID: cfg.ClientID,
	}

	return r.clientIDs[cfg.ClientID].RoundTripper, nil
}

func (r *Registry) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(DefaultTickerDuration)

		for {
			select {
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				r.clearTimer()
			}
		}
	}()
}

func (r *Registry) removeClientID(clientID string) {
	if _, ok := r.clientIDs[clientID]; !ok {
		return
	}

	r.clientIDs[clientID].Cancel()
	delete(r.clientIDs, clientID)
}

func (r *Registry) clearTimer() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for _, timer := range r.timers {
		// log.Debug().Msgf("checking timer for client id: %s", timer.ClientID)
		select {
		case <-timer.Timer.C:
			// log.Debug().Msgf("clear timer for client id: %s", timer.ClientID)
			r.removeClientID(timer.ClientID)
			if !timer.Timer.Stop() {
				select {
				case <-timer.Timer.C:
				default:
				}
			}

			timer.Timer.C = nil

			delete(r.timers, timer.ClientID)
		default:
		}
	}
}

func (r *Registry) resetTimer(clientID string) {
	if _, ok := r.timers[clientID]; !ok {
		return
	}

	if !r.timers[clientID].Timer.Stop() {
		<-r.timers[clientID].Timer.C
	}

	r.timers[clientID].Timer.Reset(DefaultTimeDuration)
}
