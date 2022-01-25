package store

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
)

func AutoMigrate(ctx context.Context, dbConn *gorm.DB) error {
	if err := dbConn.AutoMigrate(Models...); err != nil {
		return err
	}

	// first user initialize if not exist
	userOne := models.User{}
	result := dbConn.WithContext(ctx).Limit(1).Find(&userOne)

	if result.Error == nil || userOne.ID.ID != uuid.Nil {
		return nil
	}

	var err error

	addAdmin := models.User{}

	addAdmin.Admin = true

	addAdmin.ID.ID, err = uuid.NewUUID()
	if err != nil {
		return err
	}

	addAdmin.Name = config.Application.Name

	addAdmin.Password, err = sec.HashPassword(config.Application.Password)
	if err != nil {
		return err
	}

	result = dbConn.WithContext(ctx).Create(&addAdmin)

	return result.Error
}
