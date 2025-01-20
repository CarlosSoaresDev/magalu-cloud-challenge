package provider

type ProviderType string

const (
	PayPalGateway ProviderType = "PayPal"
	StripeGateway ProviderType = "Stripe"
)
