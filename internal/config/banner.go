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
╚%s %s ==╝
version=[%s] buildCommit=[%s] buildDate=[%s]
`

func Banner(txt string) string {
	//nolint:gomnd // banner static char count
	repeat := 34 - (len(txt))
	if 0 > repeat {
		repeat = 0
	}

	return fmt.Sprintf(bannerLogo, strings.Repeat("=", repeat), txt, AppVersion, AppBuildCommit, AppBuildDate)
}
