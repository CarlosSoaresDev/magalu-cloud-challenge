package processor

import (
	"errors"

	paypalService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/paypal"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/paypal/processor/actions"
)

type PaymentGateway interface {
	Process(service paypalService.PaypalService, event interface{}) error
}

var paymentGateways = map[PayPalProcessType]PaymentGateway{
	createdAction: &actions.PayPalCreatedAction{},
}

// NewProcessor creates a new instance of a PaymentGateway based on the provided PayPalProcessType.
// It returns the corresponding PaymentGateway if the processor type exists, otherwise it returns an error.
//
// Parameters:
//   - processorType: The type of PayPal process to create a PaymentGateway for.
//
// Returns:
//   - PaymentGateway: The created PaymentGateway instance.
//   - error: An error if the processor type is unsupported.
func NewProcessor(processorType PayPalProcessType) (PaymentGateway, error) {
	if gateway, exists := paymentGateways[processorType]; exists {
		return gateway, nil
	}
	return nil, errors.New("unsupported paypal action")
}
