package actions

import (
	"encoding/json"

	stripeService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe"
	"github.com/stripe/stripe-go"
)

type StripeCreatedAction struct{}

// Process handles the "payment_created" event from Stripe.
// It unmarshals the event data into a PaymentIntent object and
// adds a transaction with the status "created" using the provided StripeService.
//
// Parameters:
// - service: An instance of StripeService to interact with Stripe API.
// - event: The Stripe event containing the payment intent data.
//
// Returns:
// - error: An error if the unmarshalling or adding transaction fails, otherwise nil.
func (pg *StripeCreatedAction) Process(service stripeService.StripeService, event stripe.Event) error {

	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)

	if err != nil {
		return err
	}

	return service.AddTransaction(paymentIntent.ID, "created")
}
