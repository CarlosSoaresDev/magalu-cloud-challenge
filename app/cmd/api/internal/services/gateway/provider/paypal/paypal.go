package paypal

import (
	"fmt"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
)

type PayPalGateway struct{}

// ProcessPayment processes a payment using the PayPal gateway.
// It takes a payment of type models.Gateway and a correlationId string as parameters.
// It returns a pointer to a string and an error.
// If the payment cannot be processed, it returns an error indicating the failure reason.
func (pg *PayPalGateway) ProcessPayment(payment models.Gateway, correlationId string) (*string, error) {
	fmt.Printf("Processing PayPal payment of $%.2f\n", payment.Amount)

	return nil, fmt.Errorf("unable to process payment using gateway: %s", payment.Gateway)
}
