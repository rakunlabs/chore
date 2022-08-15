package store

import "github.com/worldline-go/chore/models"

var Models = []interface{}{
	&models.Auth{},
	&models.Template{},
	&models.Folder{},
	&models.Group{},
	&models.User{},
	&models.Token{},
	&models.Control{},
	&models.Settings{},
	// &models.Test{},
}
