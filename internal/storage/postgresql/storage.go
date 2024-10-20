package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
	"url-shortener/internal/Logger"
	"url-shortener/internal/config"
	"url-shortener/internal/entities"
	"url-shortener/pkg/random"
)

type Storage struct {
	db     *sqlx.DB
	logger Logger.Logger
}

func NewStorage(pgConfig config.PostgresServer, logger Logger.Logger) (*Storage, error) {
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

func (s *Storage) SelectAll() ([]string, error) {
	rows, err := s.db.Query(`SELECT * FROM urls`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0)
	for rows.Next() {
		var id int64
		var url string
		var alias string
		err = rows.Scan(&id, &alias, &url)
		if err != nil {
			return nil, err
		}
		res = append(res, strconv.Itoa(int(id))+" "+url+" "+alias)
	}

	s.logger.Debug("Returning from SelectAll")

	return res, nil
}

func (s *Storage) GetUniqueFreeAlias(aliasLen int) (string, error) {
	randString := random.RandomString(aliasLen)
	for {
		urls := make([]entities.URL, 0)
		err := s.db.Select(&urls, "SELECT * FROM urls WHERE alias=$1", randString)
		if err != nil {
			s.logger.Error("Error with select", err)
			return "", err
		}

		if len(urls) == 0 {
			break
		}

		randString = random.RandomString(aliasLen)
	}

	return randString, nil
}

func (s *Storage) AddUrl(url string, alias string) (*entities.URL, error) {
	var id uint64
	err := s.db.Get(&id, "INSERT INTO urls (url, alias) VALUES($1, $2) RETURNING id", url, alias)

	if err != nil {
		s.logger.Error("Error with insert", err)
		return nil, err
	}

	return &entities.URL{Id: id, Url: url, Alias: alias}, nil
}

func (s *Storage) GetUrlByAlias(alias string) (*entities.URL, error) {
	urls := make([]entities.URL, 0)
	err := s.db.Select(&urls, "SELECT * FROM urls WHERE alias=$1", alias)
	if err != nil {
		s.logger.Error("Error with select", err)
		return nil, err
	}

	if len(urls) == 0 {
		return nil, nil
	}

	return &urls[0], nil
}

func (s *Storage) DeleteUrlByAlias(alias string) (bool, error) {
	result, err := s.db.Exec("DELETE FROM urls WHERE alias=$1 RETURNING url", alias)
	if err != nil {
		s.logger.Error("Error with delete", err)
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("Error with rowsAffected", err)
		return false, err
	}

	if rowsAffected == 0 {
		s.logger.Debug("Zero rows deleted")
		return false, nil
	}

	return true, nil
}

func (s *Storage) Close() {
	s.db.Close()
}
