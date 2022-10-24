package registry

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/rytsh/liz/utils/templatex"
	"github.com/worldline-go/chore/pkg/sec"
)

type AppStore struct {
	Template *templatex.Template
	App      *fiber.App
	DB       *gorm.DB
	JWT      *sec.JWT
}

type Registry struct {
	apps  map[string]*AppStore
	WG    *sync.WaitGroup
	mutex sync.RWMutex
}

func (r *Registry) Get(name string) *AppStore {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.apps[name]
}

func (r *Registry) Iter(fn func(*AppStore)) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for k := range r.apps {
		fn(r.apps[k])
	}
}

func (r *Registry) Set(name string, appStore *AppStore) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.apps[name] = appStore
}

var (
	regOnce  sync.Once
	registry *Registry
)

// Reg get the registry or create it if not exists.
// Options can be used just once.
func Reg(options ...Option) *Registry {
	regOnce.Do(func() {
		reg := &Registry{
			apps: make(map[string]*AppStore),
		}

		for _, opt := range options {
			opt(reg)
		}

		registry = reg
	})

	return registry
}
