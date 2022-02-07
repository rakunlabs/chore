package store

import "gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"

var Models = []interface{}{
	&models.Bind{},
	&models.Auth{},
	&models.Template{},
	&models.Folder{},
	&models.User{},
	&models.Token{},
	// &models.Test{},
}
