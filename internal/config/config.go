package config

var (
	appName    = "chore"
	appVersion = "v0.0.0"
)

var Application = struct {
	Host       string
	Port       string
	LogLevel   string
	AppName    string
	StoreName  string
	AppVersion string
}{
	Host:       "0.0.0.0",
	LogLevel:   "debug",
	Port:       "8080",
	AppName:    appName,
	AppVersion: appVersion,
}
