package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"url-shortener/internal/config"
)

type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type Storage struct {
	db     *sqlx.DB
	logger Logger
}

func NewStorage(pgConfig config.PostgresServer, logger Logger) (*Storage, error) {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		pgConfig.User, pgConfig.Password, pgConfig.Address, pgConfig.Port, pgConfig.DbName)

	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		logger.Error("Error connecting to database")
		return nil, err
	}

	logger.Debug("Connected to database")

	return &Storage{conn, logger}, nil
}
