package store

import "gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"

var Models = []interface{}{
	&models.Auth{},
	&models.Template{},
	&models.Folder{},
	&models.Group{},
	&models.User{},
	&models.Token{},
	&models.Control{},
	// &models.Test{},
}
