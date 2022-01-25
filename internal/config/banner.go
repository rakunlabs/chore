package config

import (
	"fmt"
	"strings"
)

var bannerLogo = `
 █████╗ ██╗  ██╗ █████╗ ██████╗ ███████╗
██╔══██╗██║  ██║██╔══██╗██╔══██╗██╔════╝
██║  ╚═╝███████║██║  ██║██████╔╝█████╗  
██║  ██╗██╔══██║██║  ██║██╔══██╗██╔══╝  
║█████╔╝██║  ██║╚█████╔╝██║  ██║███████╗
╚%s%s [%s]==╝
`

func Banner(txt string) string {
	//nolint:gomnd // banner static char count
	return fmt.Sprintf(bannerLogo, strings.Repeat("=", 33-(len(txt)+len(AppVersion))), txt, AppVersion)
}
