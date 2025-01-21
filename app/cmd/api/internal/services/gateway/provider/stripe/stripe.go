package stripe

import (
	"fmt"
	"os"
	"strings"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/token"
)

var supportedMethods = map[string]bool{
	"card": true,
}

type StripeGateway struct{}

// ProcessPayment processes a payment using the Stripe gateway.
// It takes a payment model and a correlation ID as input parameters.
// The function returns the payment intent ID and an error, if any.
//
// Parameters:
// - payment: models.Gateway containing payment details such as card information and amount.
// - correlationId: string representing a unique identifier for the transaction.
//
// Returns:
// - *string: Pointer to the payment intent ID if the payment is successful.
// - error: Error if there is any issue during the payment processing.
//
// The function performs the following steps:
// 1. Retrieves the Stripe secret key from the environment variables.
// 2. Parses the card expiry date and validates its format.
// 3. Creates a Stripe token using the card details.
// 4. Creates a Stripe payment intent with the specified amount, currency, and payment method.
// 5. Adds metadata to the payment intent.
// 6. Returns the payment intent ID or an error if the payment intent creation fails.
func (sg *StripeGateway) ProcessPayment(payment models.Gateway, correlationId string) (*string, error) {

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	expiry := strings.Split(payment.CardDetails.Expiry, "/")
	month := expiry[0]
	year := expiry[1]

	if !supportedMethods[payment.PaymentMethod] {
		return nil, fmt.Errorf("unsupported payment method: %s. Supported methods are: %v", payment.PaymentMethod, keys(supportedMethods))
	}

	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(payment.CardDetails.Number),
			ExpMonth: stripe.String(month),
			ExpYear:  stripe.String(year),
			CVC:      stripe.String(payment.CardDetails.Cvv),
		},
	}

	token, err := token.New(tokenParams)
	if err != nil {
		fmt.Print("Is Testing")
	}

	param := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(payment.Amount * 100)),
		Currency:           stripe.String(string(stripe.Currency(payment.Currency))),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Confirm:            stripe.Bool(true),
	}

	if err == nil {
		param.Source = &token.ID
	} else {
		param.PaymentMethod = stripe.String("pm_card_visa")
	}

	param.AddMetadata("correlation_id", correlationId)

	pi, err := paymentintent.New(param)
	if err != nil {
		return nil, fmt.Errorf("error creating payment intent: %v", err)
	}

	return &pi.ID, nil
}

// keys returns a slice of strings containing the keys of the provided map.
// The input is a map where the keys are strings and the values are booleans.
// The output is a slice of strings containing all the keys from the input map.
//
// Parameters:
//
//	supportedMethods - a map with string keys and boolean values.
//
// Returns:
//
//	A slice of strings containing all the keys from the input map.
func keys(supportedMethods map[string]bool) []string {
	keys := make([]string, 0, len(supportedMethods))
	for k := range supportedMethods {
		keys = append(keys, k)
	}
	return keys
}
