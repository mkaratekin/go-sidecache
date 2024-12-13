package service

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const SideCacheHostUrl = "http://localhost:9191/"

type CacheService interface {
	GetCache(key string) ([]byte, error)
	InvalidateCache(key string, ttlAsSecond int) error
}

type cacheService struct {
	client *http.Client
}

func NewCacheService(client *http.Client) CacheService {
	return &cacheService{client: client}
}

func (c *cacheService) GetCache(key string) ([]byte, error) {
	url := SideCacheHostUrl + key

	// Make GET request
	resp, err := c.makeRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !isSuccess(resp.StatusCode) {
		return nil, fmt.Errorf("cache returns unsuccessful response code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// InvalidateCache invalidates the cache for a given key with a TTL.
func (c *cacheService) InvalidateCache(key string, ttlAsSecond int) error {
	url := SideCacheHostUrl + key

	ttlHeaderValue := "ttl=" + strconv.Itoa(ttlAsSecond)
	headers := map[string]string{
		"Content-Type":      "application/json",
		"tysidecarcachable": ttlHeaderValue,
	}

	// Make PUT request
	resp, err := c.makeRequest(http.MethodPut, url, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *cacheService) makeRequest(method, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set request headers if provided
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 400
}
