package provider

import (
	"errors"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/gateway/provider/paypal"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/gateway/provider/stripe"
)

type PaymentGateway interface {
	ProcessPayment(payment models.Gateway, correlationId string) (*string, error)
}

var Providers = map[ProviderType]PaymentGateway{
	PayPalGateway: &paypal.PayPalGateway{},
	StripeGateway: &stripe.StripeGateway{},
}

// NewProvider creates a new instance of a PaymentGateway based on the provided ProviderType.
// It returns the corresponding PaymentGateway if the type exists in the Providers map,
// otherwise, it returns an error indicating that the payment gateway type is unsupported.
//
// Parameters:
//   - gwType: The type of the payment gateway to be created.
//
// Returns:
//   - PaymentGateway: The created payment gateway instance.
//   - error: An error if the payment gateway type is unsupported.
func NewProvider(gwType ProviderType) (PaymentGateway, error) {
	if gateway, exists := Providers[gwType]; exists {
		return gateway, nil
	}
	return nil, errors.New("unsupported payment gateway type")
}
