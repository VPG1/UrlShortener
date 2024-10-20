package services

import (
	"url-shortener/internal/Logger"
	"url-shortener/internal/entities"
)

type Storage interface {
	SelectAll() ([]string, error)
	AddUrl(string, string) (*entities.URL, error)
	GetUrlByAlias(string) (*entities.URL, error)
	GetUniqueFreeAlias(int) (string, error)
	DeleteUrlByAlias(string) (bool, error)
}

type UrlService struct {
	AliasLen int
	Storage  Storage
	Logger   Logger.Logger
}

func NewUrlService(aliasLen int, storage Storage, logger Logger.Logger) *UrlService {
	return &UrlService{AliasLen: aliasLen, Storage: storage, Logger: logger}
}

func (us *UrlService) GetUrls() ([]string, error) {
	us.Logger.Debug("Getting urls")

	urls, err := us.Storage.SelectAll()
	if err != nil {
		us.Logger.Error("Failed to get urls", "error", err)
		return nil, err
	}

	return urls, nil
}

func (us *UrlService) CreateNewAlias(url string) (*entities.URL, error) {
	us.Logger.Debug("Creating new alias")

	alias, err := us.Storage.GetUniqueFreeAlias(us.AliasLen)
	if err != nil {
		us.Logger.Error("Failed to get free alias", "url", url, "error", err)
		return nil, err
	}

	urlEntity, err := us.Storage.AddUrl(url, alias)
	if err != nil {
		us.Logger.Error("Failed to add url", "url", url, "err", err)
		return nil, err
	}

	us.Logger.Info("Alias created", "url", url, "alias", alias)

	return urlEntity, nil
}

func (us *UrlService) GetUrlByAlias(alias string) (*entities.URL, error) {
	us.Logger.Debug("Getting url by alias", "alias", alias)

	url, err := us.Storage.GetUrlByAlias(alias)
	if err != nil {
		us.Logger.Error("Failed to get url by alias", "url", url, "err", err)
	}

	if url == nil {
		us.Logger.Info("Url not found by alias", "alias", alias)

		return url, nil
	}

	us.Logger.Info("Url found by alias", "url", url)

	return url, nil
}

func (us *UrlService) DeleteUrlByAlias(alias string) (bool, error) {
	us.Logger.Debug("Deleting url", "alias", alias)
	isUrlDeleted, err := us.Storage.DeleteUrlByAlias(alias)
	if err != nil {
		us.Logger.Error("Failed to delete url by alias", "err", err)
		return false, err
	}

	return isUrlDeleted, nil
}
