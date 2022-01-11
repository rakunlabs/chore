package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type DB uint

const (
	TypePostgres DB = iota
)

type ConnConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func OpenConnection(typeDB DB, cfg *ConnConfig) *gorm.DB {
	db, err := PostgresDB(cfg)
	if err != nil {
		log.Error().Err(err).Msg("cannot open connection to db")

		return nil
	}

	return db
}
