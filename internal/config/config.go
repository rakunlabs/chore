package config

import "time"

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
	Env      string `cfg:"env"`
	Host     string `cfg:"host"`
	Port     string `cfg:"port"`
	LogLevel string `cfg:"logLevel"`
	Secret   string `cfg:"secret" loggable:"false"`
	// BasePath for swagger ui ex /chore/
	BasePath string `cfg:"basePath"`
	User     User   `cfg:"user"`
	Store    Store  `cfg:"store"`
	Migrate  Store  `cfg:"migrate"`
	Server   Server `cfg:"server"`
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
		ConnMaxLifetime: 15 * time.Minute,
		MaxIdleConns:    5,
		MaxOpenConns:    7,
	},
	Server: Server{
		ReadBufferSize:  1024 * 1024,
		WriteBufferSize: 1024 * 1024,
	},
}

// User settings will use if doesn't have any user on database.
type User struct {
	Name     string `cfg:"name"`
	Password string `cfg:"password" loggable:"false"`
}

type Store struct {
	Type            string        `cfg:"type"`
	FileName        string        `cfg:"fileName"`
	Schema          string        `cfg:"schema"`
	Host            string        `cfg:"host"`
	Port            string        `cfg:"port"`
	User            string        `cfg:"user"`
	Password        string        `cfg:"password" loggable:"false"`
	DBName          string        `cfg:"dbName"`
	TimeZone        string        `cfg:"timeZone" default:"UTC"`
	DBDataSource    string        `cfg:"dbDataSource" loggable:"false"`
	ConnMaxIdleTime time.Duration `cfg:"connMaxIdleTime"`
	ConnMaxLifetime time.Duration `cfg:"connMaxLifetime"`
	MaxIdleConns    int           `cfg:"maxIdleConns"`
	MaxOpenConns    int           `cfg:"maxOpenConns"`
}

type Server struct {
	ReadBufferSize  int `cfg:"readBufferSize"`
	WriteBufferSize int `cfg:"writeBufferSize"`
}
