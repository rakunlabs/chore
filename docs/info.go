package docs

import (
	"github.com/worldline-go/swagger"
)

func SetInfo(title, version, basePath string) error {
	options := []swagger.Option{swagger.WithTitle(title), swagger.WithVersion(version)}
	if basePath != "" {
		options = append(options, swagger.WithBasePath(basePath))
	}

	return swagger.SetInfo(options...) //nolint:wrapcheck // no need
}
