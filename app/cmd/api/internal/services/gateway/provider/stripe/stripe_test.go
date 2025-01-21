package stripe

import (
	"os"
	"testing"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func setupMockEnvironment() {
	os.Setenv("STRIPE_SECRET_KEY", "sk_test_4eC39HqLyjWDarjtT1zdp7dc")
}

func TestProcessPayment_Successful(t *testing.T) {
	// Arrange
	sg := &StripeGateway{}
	payment := models.Gateway{
		Gateway:       "Stripe",
		Amount:        100.00,
		Currency:      "USD",
		PaymentMethod: "card",
		CardDetails: models.CardDetails{
			Expiry: "12/23",
		},
	}
	setupMockEnvironment()

	// Action
	correlationId := utils.GenerateGUID()
	_, err := sg.ProcessPayment(payment, correlationId)

	// Assert
	assert.NoError(t, err)
}

func TestProcessPayment_InvalidPaymentMethod(t *testing.T) {
	// Arrange
	sg := &StripeGateway{}
	payment := models.Gateway{}
	payment.PaymentMethod = "1234"
	payment.CardDetails.Expiry = "12/23"
	setupMockEnvironment()

	// Action
	correlationId := utils.GenerateGUID()
	paymentIntentID, err := sg.ProcessPayment(payment, correlationId)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, paymentIntentID)
	assert.Equal(t, "unsupported payment method: 1234. Supported methods are: [card]", err.Error())
}
