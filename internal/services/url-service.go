package services

import (
	"url-shortener/internal/entities"
)

type Storage interface {
	AddUrl(string, string) (*entities.URL, error)
	GetUrlByAlias(string) (*entities.URL, error)
	GetUniqueFreeAlias(int) (string string, err error)
}

type Logger interface {
	Info(string, ...any)
	Debug(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type UrlService struct {
	AliasLen int
	Storage  *Storage
	Logger   *Logger
}

func NewUrlService(aliasLen int, storage *Storage, logger *Logger) *UrlService {
	return &UrlService{AliasLen: aliasLen, Storage: storage, Logger: logger}
}

func (us *UrlService) CreateNewAlias(url string) (*entities.URL, error) {
	(*us.Logger).Debug("Creating new alias")

	alias, err := (*us.Storage).GetUniqueFreeAlias(us.AliasLen)
	if err != nil {
		(*us.Logger).Error("Failed to get free alias", "url", url, "error", err)
		return nil, err
	}

	urlEntity, err := (*us.Storage).AddUrl(url, alias)
	if err != nil {
		(*us.Logger).Error("Failed to add url", "url", url, "err", err)
		return nil, err
	}

	(*us.Logger).Info("Alias created", "url", url, "alias", alias)

	return urlEntity, nil
}

func (us *UrlService) GetUrlByAlias(alias string) (*entities.URL, error) {
	url, err := (*us.Storage).GetUrlByAlias(alias)
	if err != nil {
		(*us.Logger).Error("Failed to get url by alias", "url", url, "err", err)
	}

	(*us.Logger).Info("Alias found", "url", url)

	return url, nil
}
