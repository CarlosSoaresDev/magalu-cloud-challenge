package paypal

import (
	"testing"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
)

func TestProcessPayment(t *testing.T) {
	tests := []struct {
		name          string
		payment       models.Gateway
		correlationId string
		wantErr       bool
	}{
		{
			name: "valid payment",
			payment: models.Gateway{
				Amount:  100.00,
				Gateway: "PayPal",
			},
			correlationId: "12345",
			wantErr:       true,
		},
		{
			name: "invalid payment gateway",
			payment: models.Gateway{
				Amount:  50.00,
				Gateway: "InvalidGateway",
			},
			correlationId: "67890",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PayPalGateway{}
			got, err := pg.ProcessPayment(tt.payment, tt.correlationId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessPayment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Errorf("ProcessPayment() got = %v, want nil", *got)
			}
		})
	}
}
