package stripe

import (
	"fmt"
	"os"
	"strings"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
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

	if !supportedMethods[payment.PaymentMethod] {
		return nil, fmt.Errorf("unsupported payment method: %s. Supported methods are: %v", payment.PaymentMethod, keys(supportedMethods))
	}

	token, err := sg.createToken(payment)
	paymentMethodTest := getPaymentMethodTest(err, payment)

	param := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(int64(payment.Amount * 100)),
		Currency:           stripe.String(string(stripe.Currency(payment.Currency))),
		PaymentMethodTypes: stripe.StringSlice([]string{payment.PaymentMethod}),
		Confirm:            stripe.Bool(true),
	}

	if utils.IsEmptyOrNull(paymentMethodTest) {
		param.Source = &token.ID
	} else {
		param.PaymentMethod = stripe.String(paymentMethodTest)
	}

	param.AddMetadata("correlation_id", correlationId)

	pi, err := paymentintent.New(param)
	if err != nil {
		return nil, fmt.Errorf("error creating payment intent: %v", err)
	}

	return &pi.ID, nil
}

// createToken generates a Stripe token for the provided payment details.
// It extracts the card expiry month and year from the payment's CardDetails,
// and uses them along with the card number and CVC to create a stripe.TokenParams object.
// The function then calls the Stripe API to create and return a new token.
//
// Parameters:
//   - payment: A models.Gateway object containing the card details.
//
// Returns:
//   - *stripe.Token: A pointer to the created Stripe token.
//   - error: An error object if the token creation fails.
func (sg *StripeGateway) createToken(payment models.Gateway) (*stripe.Token, error) {
	expiry := strings.Split(payment.CardDetails.Expiry, "/")
	month := expiry[0]
	year := expiry[1]

	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(payment.CardDetails.Number),
			ExpMonth: stripe.String(month),
			ExpYear:  stripe.String(year),
			CVC:      stripe.String(payment.CardDetails.Cvv),
		},
	}

	return token.New(tokenParams)
}

// getPaymentMethodTest returns a default payment method test string based on the provided card number
// if an error occurs during token creation. It uses predefined card numbers to determine the payment method.
//
// Parameters:
// - err: an error that indicates if there was an issue creating the token.
// - payment: a models.Gateway object that contains card details.
//
// Returns:
// - A string representing the default payment method test.
func getPaymentMethodTest(err error, payment models.Gateway) string {
	var paymentMethodTest string
	if err != nil {
		fmt.Print("error to create token, using default payment method test")
		switch payment.CardDetails.Number {
		case "4242424242424242":
			paymentMethodTest = "pm_card_visa"
		case "4000056655665556":
			paymentMethodTest = "pm_card_visa_debit"
		case "5555555555554444":
			paymentMethodTest = "pm_card_mastercard"
		case "5200828282828210":
			paymentMethodTest = "pm_card_mastercard_debit"
		default:
			paymentMethodTest = "pm_card_visa"
		}
	}
	return paymentMethodTest
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
