package config

import "time"

var (
	AppName    = "chore"
	AppVersion = "v0.0.0"
	LoadName   = ""
	Env        = "DEVELOPMENT"
	StartDate  = time.Now()
)

var Application = struct {
	Host     string `cfg:"host"`
	Port     string `cfg:"port"`
	LogLevel string `cfg:"logLevel"`
	Secret   string `cfg:"secret" loggable:"false"`
	// BasePath for swagger ui ex /chore/
	BasePath string `cfg:"basePath"`
	User     User   `cfg:"user"`
	Store    Store  `cfg:"store"`
	Server   Server `cfg:"server"`
}{
	Host:     "0.0.0.0",
	LogLevel: "debug",
	Port:     "8080",
	User: User{
		Name:     "admin",
		Password: "admin",
	},
	Store: Store{
		Schema: AppName,
		Type:   "postgres",
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
	Type     string `cfg:"type"`
	FileName string `cfg:"fileName"`
	Schema   string `cfg:"schema"`
	Host     string `cfg:"host"`
	Port     string `cfg:"port"`
	User     string `cfg:"user"`
	Password string `cfg:"password" loggable:"false"`
	DBName   string `cfg:"dbName"`
}

type Server struct {
	ReadBufferSize  int `cfg:"readBufferSize"`
	WriteBufferSize int `cfg:"writeBufferSize"`
}
