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

func NewProcessor(processorType StripeProcessType) (StripeProcessor, error) {
	if gateway, exists := paymentGateways[processorType]; exists {
		return gateway, nil
	}
	return nil, errors.New("unsupported stripe action")
}
