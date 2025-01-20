package gateway

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/gateway/provider"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"
)

type GatewayService interface {
	GetAllAvaiablesGateways() []string
	GetAllTransactionsByDate(date string) (*[]models.Transaction, error)
	AddTransaction(id string, payment models.Gateway) error
}

type gatewayService struct {
	cache cache.CacheClient
}

// New creates a new instance of gatewayService with the provided cache client.
// It returns a pointer to the newly created gatewayService.
//
// Parameters:
//   - cache: an instance of cache.CacheClient to be used by the gatewayService.
//
// Returns:
//   - *gatewayService: a pointer to the newly created gatewayService.
func New(cache cache.CacheClient) *gatewayService {
	return &gatewayService{
		cache: cache,
	}
}

// GetAllAvaiablesGateways retrieves all available gateway keys from the provider.
// It returns a slice of strings containing the keys of the available gateways.
func (p *gatewayService) GetAllAvaiablesGateways() []string {

	keys := make([]string, 0, len(provider.Providers))

	for key := range provider.Providers {
		keys = append(keys, string(key))
	}

	return keys
}

// GetAllTransactionsByDate retrieves all transactions for a given date from the cache.
// The date parameter should be in the format "YYYY/MM/DD".
// It returns a pointer to a slice of Transaction models and an error if any occurs during the process.
// If the transactions are not found in the cache, it returns an empty slice and a nil error.
//
// Parameters:
//   - date: A string representing the date for which transactions are to be retrieved.
//
// Returns:
//   - *[]models.Transaction: A pointer to a slice of Transaction models.
//   - error: An error if any occurs during the process, otherwise nil.
func (p *gatewayService) GetAllTransactionsByDate(date string) (*[]models.Transaction, error) {

	var transactions []models.Transaction
	var transactionsMap map[string]models.Transaction
	transactionsByDate := fmt.Sprintf("%s_%s", cache.TransactionsKey, strings.Replace(date, "/", "_", 2))
	c, err := p.cache.Get(transactionsByDate)

	if err != nil {
		return &[]models.Transaction{}, nil
	}

	err = json.Unmarshal(c, &transactionsMap)

	if err != nil {
		return nil, err
	}

	for _, transaction := range transactionsMap {
		transactions = append(transactions, transaction)
	}

	return &transactions, nil
}

// AddTransaction adds a new transaction to the cache with the given id and payment details.
// It creates a new transaction with the current timestamp and a status of "pending".
// The transaction is then stored in the cache, grouped by the current date.
//
// Parameters:
//   - id: A string representing the unique identifier for the transaction.
//   - payment: A models.Gateway object containing the payment details.
//
// Returns:
//   - error: An error if there is an issue with cache retrieval, unmarshalling, or setting the cache.
func (p *gatewayService) AddTransaction(id string, payment models.Gateway) error {
	now := time.Now()
	transaction := models.Transaction{
		Id:       id,
		Amount:   payment.Amount,
		Currency: payment.Currency,
		TransactionStatus: []models.TransactionStatus{
			{
				DateTime: now.Format(time.RFC3339),
				Status:   "pending",
			},
		},
	}

	transactionsByDate := fmt.Sprintf("%s_%s", cache.TransactionsKey, now.Format("02_01_2006"))
	c, err := p.cache.Get(transactionsByDate)

	var transactions map[string]models.Transaction
	if err != nil {
		if err.Error() == cache.ErrCacheMiss.Error() {
			transactions = map[string]models.Transaction{id: transaction}
		} else {
			return err
		}
	} else {
		if err := json.Unmarshal(c, &transactions); err != nil {
			return err
		}
		transactions[id] = transaction
	}

	transactionsSerialized, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	if err := p.cache.Set(transactionsByDate, string(transactionsSerialized), 0); err != nil {
		return err
	}

	return nil
}
