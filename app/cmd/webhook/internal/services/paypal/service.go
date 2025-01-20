package paypal

import "github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"

type PaypalService interface {
}

type paypalService struct {
	cache cache.CacheClient
}

// New creates a new instance of paypalService with the provided cache client.
//
// Parameters:
//   - cache: an instance of CacheClient to be used for caching purposes.
//
// Returns:
//   - *paypalService: a new instance of paypalService.
func New(cache cache.CacheClient) *paypalService {
	return &paypalService{
		cache: cache,
	}
}
