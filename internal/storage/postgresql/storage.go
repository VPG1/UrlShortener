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

func (s *Storage) SelectAllUserId(userId uint64) ([]string, error) {
	rows, err := s.db.Query("SELECT * FROM urls WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]string, 0)
	for rows.Next() {
		var id int64
		var url string
		var alias string
		var userId uint64
		err = rows.Scan(&id, &alias, &url, &userId)
		if err != nil {
			return nil, err
		}
		res = append(res, strconv.Itoa(int(id))+", "+url+", "+alias+", "+strconv.Itoa(int(userId)))
	}

	s.logger.Debug("Returning from SelectAllUserId")

	return res, nil
}

func (s *Storage) GetUniqueFreeAlias(aliasLen int) (string, error) {
	randString := random.RandomString(aliasLen)
	for {
		urls := make([]entities.URL, 0)
		err := s.db.Select(&urls, "SELECT * FROM urls WHERE alias=$1", randString)
		if err != nil {
			s.logger.Error("Error with select", "err", err)
			return "", err
		}

		if len(urls) == 0 {
			break
		}

		randString = random.RandomString(aliasLen)
	}

	return randString, nil
}

func (s *Storage) AddUrl(url string, alias string, userId uint64) (*entities.URL, error) {
	var id uint64
	err := s.db.Get(&id, "INSERT INTO urls (url, alias, user_id) VALUES($1, $2, $3) RETURNING id", url, alias, userId)

	if err != nil {
		s.logger.Error("Error with insert", err)
		return nil, err
	}

	return &entities.URL{Id: id, Url: url, Alias: alias, UserId: userId}, nil
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

func (s *Storage) DeleteUrlByAlias(alias string, userId uint64) (bool, error) {

	result, err := s.db.Exec("DELETE FROM urls WHERE alias=$1 AND user_id=$2 RETURNING url", alias, userId)
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

func (s *Storage) AddUser(name string, userName string, passwordHash string) (*entities.User, error) {
	var id uint64
	err := s.db.Get(&id, "INSERT INTO users (name, username, password_hash) VALUES($1, $2, $3) RETURNING id",
		name, userName, passwordHash)
	if err != nil {
		s.logger.Error("Error with insert", err)
		return nil, err
	}

	return &entities.User{Id: id, Name: name, UserName: userName, PasswordHash: passwordHash}, nil
}

func (s *Storage) GetUser(userName string, passwordHash string) (*entities.User, error) {
	var user []entities.User
	err := s.db.Select(&user, "SELECT * FROM users WHERE username=$1 AND password_hash=$2",
		userName, passwordHash)
	if err != nil {
		s.logger.Error("Error with select", err)
		return nil, err
	}

	if len(user) == 0 {
		return nil, nil
	}

	return &user[0], nil
}

//func (s *Storage) ParseToken(token string) (uint64, error) {
//	//jwt.ParseWithClaims(token, &c)
//}

func (s *Storage) Close() {
	s.db.Close()
}
