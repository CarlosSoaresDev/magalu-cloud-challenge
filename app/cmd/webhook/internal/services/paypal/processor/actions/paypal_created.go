package actions

import (
	"fmt"

	paypalService "github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/webhook/internal/services/paypal"
)

type PayPalCreatedAction struct{}

// Process processes the PayPal created event.
// It takes a PayPalService and a PayPal event as parameters and returns an error if any occurs during processing.
//
// Parameters:
// - service: an instance of PayPalService to handle PayPal related operations.
// - event: a PayPal event that contains the details of the PayPal created event.
//
// Returns:
// - error: an error if any issue occurs during the processing of the event, otherwise nil.
func (pg *PayPalCreatedAction) Process(service paypalService.PaypalService, event interface{}) error {
	fmt.Printf("Processing PayPal")

	return nil
}
