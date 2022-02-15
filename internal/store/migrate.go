package store

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/internal/config"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/models"
	"gitlab.test.igdcs.com/finops/nextgen/apps/tools/chore/pkg/sec"
)

func AutoMigrate(ctx context.Context, dbConn *gorm.DB) error {
	if err := dbConn.AutoMigrate(Models...); err != nil {
		return err
	}

	if err := initializeAdminGroup(ctx, dbConn); err != nil {
		return err
	}

	if err := initializeAdminUser(ctx, dbConn); err != nil {
		return err
	}

	return nil
}

func initializeAdminGroup(ctx context.Context, dbConn *gorm.DB) error {
	// initialize admin group
	groupAdmin := models.Group{}
	result := dbConn.WithContext(ctx).Where("name = ?", "admin").Find(&groupAdmin)

	if result.RowsAffected != 0 {
		return nil
	}

	groupAdmin.Name = "admin"
	groupAdmin.ID.ID = uuid.New()

	result = dbConn.WithContext(ctx).Create(&groupAdmin)

	return result.Error
}

func initializeAdminUser(ctx context.Context, dbConn *gorm.DB) error {
	// first admin initialize if not exist
	userOne := models.User{}
	result := dbConn.WithContext(ctx).Limit(1).Find(&userOne)

	if result.RowsAffected != 0 {
		return nil
	}

	addAdmin := models.User{}

	addAdmin.Groups.Groups = datatypes.JSON(`["admin"]`)

	addAdmin.ID.ID = uuid.New()

	addAdmin.Name = config.Application.User.Name

	hashPass, err := sec.HashPassword([]byte(config.Application.User.Password))
	if err != nil {
		return err
	}

	addAdmin.Password = string(hashPass)

	result = dbConn.WithContext(ctx).Create(&addAdmin)

	return result.Error
}
