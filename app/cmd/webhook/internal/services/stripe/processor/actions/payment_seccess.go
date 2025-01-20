package actions

import (
	"encoding/json"

	stripeService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe"
	"github.com/stripe/stripe-go"
)

type StripeSuccessAction struct{}

func (pg *StripeSuccessAction) Process(service stripeService.StripeService, event stripe.Event) error {

	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)

	if err != nil {
		return err
	}

	return service.AddTransaction(paymentIntent.ID, "success")
}
