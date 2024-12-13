package go_sidecache

import (
	"errors"
	"go-sidecache/service"
	"net/http"
)

type Sidecache interface {
	GetCache(key string) ([]byte, error)
	InvalidateCache(key string, ttl int) error
}

type sidecache struct {
	client       *http.Client
	cacheService service.CacheService
}

func NewSidecache(client *http.Client) (Sidecache, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	cacheService := service.NewCacheService(client)

	return &sidecache{client: client, cacheService: cacheService}, nil

}

func (s *sidecache) GetCache(key string) ([]byte, error) {
	return s.cacheService.GetCache(key)
}

func (s *sidecache) InvalidateCache(key string, ttl int) error {
	return s.cacheService.InvalidateCache(key, ttl)
}
