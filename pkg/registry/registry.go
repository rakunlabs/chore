package registry

import (
	"sync"

	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/rytsh/mugo/pkg/templatex"
	"github.com/worldline-go/auth"
	"github.com/worldline-go/auth/providers"
)

type Registry struct {
	Template      *templatex.Template
	Server        *echo.Echo
	DB            *gorm.DB
	JWT           JWT
	WG            *sync.WaitGroup
	AuthProviders map[string]*providers.Generic
}

type JWT struct {
	*auth.JWT
	Parser auth.JwkKeyFuncParse
}

var Reg *Registry

func Init(reg *Registry) {
	Reg = reg
}
