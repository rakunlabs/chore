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
	repeat := 33 - (len(txt) + len(AppVersion))
	if 0 > repeat {
		repeat = 0
	}

	return fmt.Sprintf(bannerLogo, strings.Repeat("=", repeat), txt, AppVersion)
}
