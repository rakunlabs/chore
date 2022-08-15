package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
	gorm_zerolog "github.com/wei840222/gorm-zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func PostgresDB(cfg map[string]interface{}) (*gorm.DB, error) {
	timeZone := cfg["timeZone"]
	if timeZone == "" {
		timeZone = "UTC"
	}

	dsn, _ := cfg["dsn"].(string)
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
			cfg["host"], cfg["user"], cfg["password"], cfg["dbName"], cfg["port"], timeZone)
	}

	gLog := gorm_zerolog.NewWithLogger(log.With().Str("component", "postgres").Logger())
	gLog.SkipErrRecordNotFound = true

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.", cfg["schema"]),
			SingularTable: false,
		},
		Logger: gLog,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres connection; %w", err)
	}

	return db, nil
}
