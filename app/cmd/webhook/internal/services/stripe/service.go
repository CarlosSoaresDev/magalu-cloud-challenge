package gateway

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"
)

type StripeService interface {
	AddTransaction(id string, status string) error
}

type stripeService struct {
	cache cache.CacheClient
}

// New creates a new instance of stripeService with the provided cache client.
// It returns a pointer to the newly created stripeService.
//
// Parameters:
//   - cache: an instance of cache.CacheClient to be used by the stripeService.
//
// Returns:
//   - *stripeService: a pointer to the newly created stripeService.
func New(cache cache.CacheClient) *stripeService {
	return &stripeService{
		cache: cache,
	}
}

// AddTransaction adds a new transaction status to an existing transaction in the cache.
// It retries up to 3 times if there are issues with unmarshalling or setting the cache.
//
// Parameters:
//   - id: The unique identifier of the transaction.
//   - status: The status to be added to the transaction.
//
// Returns:
//   - error: An error if there is an issue with unmarshalling the cache data or setting the updated transactions in the cache.
func (p *stripeService) AddTransaction(id string, status string) error {
	now := time.Now()

	var transactions = make(map[string]*models.Transaction)
	transactionsByDate := fmt.Sprintf("%s_%s", cache.TransactionsKey, now.Format("02_01_2006"))

	const maxRetries = 3
	const sleepTime = 2 * time.Second

	for i := 0; i < maxRetries; i++ {

		c, _ := p.cache.Get(transactionsByDate)

		if err := json.Unmarshal(c, &transactions); err != nil {
			return err
		}

		transaction := transactions[id]

		if transaction != nil {
			transaction.TransactionStatus = append(transaction.TransactionStatus, models.TransactionStatus{
				Status:   status,
				DateTime: now.Format(time.RFC3339),
			})

			transactions[id] = transaction

			updatedTransactions, err := json.Marshal(transactions)
			if err != nil {
				return err
			}

			if err := p.cache.Set(transactionsByDate, updatedTransactions, 0); err != nil {
				return err
			}
			break
		}
		time.Sleep(sleepTime)
	}

	return nil
}
