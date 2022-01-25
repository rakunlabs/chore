package registry

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/request"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/translate"
)

type AppStore struct {
	Template *translate.Template
	Client   *request.Client
	App      *fiber.App
	DB       *gorm.DB
	JWT      *sec.JWT
}

type Registry struct {
	apps  map[string]*AppStore
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

func Reg() *Registry {
	regOnce.Do(func() {
		registry = &Registry{
			apps: make(map[string]*AppStore),
		}
	})

	return registry
}
