package config

import (
	"time"

	"github.com/worldline-go/auth/providers"
	"github.com/worldline-go/tell"
)

var (
	AppName        = "chore"
	AppVersion     = "v0.0.0"
	AppBuildCommit = "0000000"
	AppBuildDate   = "1970-01-01T00:00:00Z"
	LoadName       = ""
	StartDate      = time.Now()
)

type Prefix struct {
	Vault  string `cfg:"vault"`
	Consul string `cfg:"consul"`
}

type Load struct {
	Prefix  Prefix `cfg:"prefix"`
	AppName string `cfg:"app_name"`
}

var LoadConfig = Load{
	AppName: AppName,
}

var Application = struct {
	Env      string   `cfg:"env"`
	Host     string   `cfg:"host"`
	Port     string   `cfg:"port"`
	LogLevel string   `cfg:"log_level"`
	Secret   string   `cfg:"secret"    log:"false"`
	BasePath string   `cfg:"base_path"`
	User     User     `cfg:"user"`
	Store    Store    `cfg:"store"`
	Migrate  Store    `cfg:"migrate"`
	Template Template `cfg:"template"`

	AuthProviders map[string]*providers.Generic `cfg:"auth_providers"`

	Telemetry tell.Config
}{
	Host:     "0.0.0.0",
	LogLevel: "info",
	Port:     "8080",
	User: User{
		Name:     "admin",
		Password: "admin",
	},
	Store: Store{
		Schema:          AppName,
		Type:            "postgres",
		User:            "postgres",
		DBName:          "postgres",
		Host:            "127.0.0.1",
		Port:            "5432",
		ConnMaxLifeTime: 15 * time.Minute,
		MaxIdleConns:    5,
		MaxOpenConns:    7,
	},
	Template: Template{
		Trust: false,
	},
}

// User settings will use if doesn't have any user on database.
type User struct {
	Name     string `cfg:"name"`
	Password string `cfg:"password" log:"false"`
}

type Store struct {
	Type            string        `cfg:"type"`
	FileName        string        `cfg:"file_name"`
	Schema          string        `cfg:"schema"`
	Host            string        `cfg:"host"`
	Port            string        `cfg:"port"`
	User            string        `cfg:"user"`
	Password        string        `cfg:"password"           log:"false"`
	DBName          string        `cfg:"db_name"`
	TimeZone        string        `cfg:"time_zone"          default:"UTC"`
	DBDataSource    string        `cfg:"db_data_source"     log:"false"`
	ConnMaxIdleTime time.Duration `cfg:"conn_max_idle_time"`
	ConnMaxLifeTime time.Duration `cfg:"conn_max_life_time"`
	MaxIdleConns    int           `cfg:"max_idle_conns"`
	MaxOpenConns    int           `cfg:"max_open_conns"`
}

type Template struct {
	Trust bool `cfg:"trust"`
}
