package config

var (
	AppName    = "chore"
	AppVersion = "v0.0.0"
	LoadName   = ""
)

var Application = struct {
	Host     string `cfg:"host" secret:"host,loggable"`
	Port     string `cfg:"port" secret:"port,loggable"`
	LogLevel string `cfg:"logLevel" secret:"logLevel,loggable"`
	Secret   string `cfg:"secret" secret:"secret"`
	Name     string `cfg:"name" secret:"name,loggable"`
	Password string `cfg:"password" secret:"password"`
	Store    Store  `cfg:"store" secret:"store,loggable"`
}{
	Host:     "0.0.0.0",
	LogLevel: "debug",
	Port:     "8080",
	Name:     "admin",
	Password: "admin",
	Store: Store{
		Schema: AppName,
		Type:   "postgres",
	},
}

type Store struct {
	Type     string `cfg:"type" secret:"type,loggable"`
	FileName string `cfg:"fileName" secret:"fileName,loggable"`
	Schema   string `cfg:"schema" secret:"schema,loggable"`
	Host     string `cfg:"host" secret:"host,loggable"`
	Port     string `cfg:"port" secret:"port,loggable"`
	User     string `cfg:"user" secret:"user,loggable"`
	Password string `cfg:"password" secret:"password"`
	DBName   string `cfg:"dbName" secret:"dbName,loggable"`
}
