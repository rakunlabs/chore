package docs

import (
	"github.com/worldline-go/swagger"
)

func SetInfo(title, version string) error {
	return swagger.SetInfo( //nolint:wrapcheck // no need
		swagger.WithTitle(title),
		swagger.WithVersion(version),
	)
}
