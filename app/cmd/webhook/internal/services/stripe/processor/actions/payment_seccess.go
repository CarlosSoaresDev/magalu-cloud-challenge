package actions

import (
	"encoding/json"

	stripeService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/stripe"
	"github.com/stripe/stripe-go"
)

type StripeSuccessAction struct{}

// Process handles the processing of a successful Stripe payment event.
// It unmarshals the event data into a PaymentIntent object and adds a transaction
// with the status "success" using the provided StripeService.
//
// Parameters:
// - service: An instance of StripeService used to add the transaction.
// - event: The Stripe event containing the payment data.
//
// Returns:
// - error: An error if the unmarshalling of event data fails or if adding the transaction fails.
func (pg *StripeSuccessAction) Process(service stripeService.StripeService, event stripe.Event) error {

	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)

	if err != nil {
		return err
	}

	return service.AddTransaction(paymentIntent.ID, "success")
}
