package processor

type StripeProcessType string

const (
	createdAction StripeProcessType = "payment_intent.created"
	successAction StripeProcessType = "payment_intent.succeeded"
)
