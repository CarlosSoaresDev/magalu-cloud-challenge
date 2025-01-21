package processor

import (
	"errors"

	stripeService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe/processor/actions"
	"github.com/stripe/stripe-go"
)

type StripeProcessor interface {
	Process(service stripeService.StripeService, event stripe.Event) error
}

var paymentGateways = map[StripeProcessType]StripeProcessor{
	createdAction: &actions.StripeCreatedAction{},
	successAction: &actions.StripeSuccessAction{},
}

// NewProcessor creates a new StripeProcessor based on the provided StripeProcessType.
// It returns the corresponding StripeProcessor if the processor type exists in the paymentGateways map.
// If the processor type does not exist, it returns an error indicating that the stripe action is unsupported.
//
// Parameters:
//   - processorType: The type of Stripe process to be created.
//
// Returns:
//   - StripeProcessor: The created StripeProcessor if the processor type exists.
//   - error: An error if the processor type is unsupported.
func NewProcessor(processorType StripeProcessType) (StripeProcessor, error) {
	if gateway, exists := paymentGateways[processorType]; exists {
		return gateway, nil
	}
	return nil, errors.New("unsupported stripe action")
}
