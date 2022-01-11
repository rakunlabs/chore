package config

var bannerLogo = `
 █████  ██   ██  █████  ██████  ███████
██   ██ ██   ██ ██   ██ ██   ██ ██     
██      ███████ ██   ██ ██████  █████  
██   ██ ██   ██ ██   ██ ██   ██ ██     
 █████  ██   ██  █████  ██   ██ ███████
`

func Banner(txt string) string {
	return bannerLogo + txt + " [" + Application.AppVersion + "]"
}
