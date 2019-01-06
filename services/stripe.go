package services

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
	"math"
)

func StripeCharge(token string, des string, amount float64) (*stripe.Charge, error) {
	total := int64(math.Round(amount* 100))

	chargeParams := &stripe.ChargeParams{
		Amount: stripe.Int64(total),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(des),
	}
	err := chargeParams.SetSource(token)

	if err != nil {
		return nil, err
	}

	ch, err := charge.New(chargeParams)

	return ch, err
}

func StripeRefund(token string) (interface{}, error) {
	refundParams := &stripe.RefundParams{
		Charge: stripe.String(token),
	}

	ref, err := refund.New(refundParams)
	return ref, err
}

func StripePartialRefund(token string, amount float64) (interface{}, error) {
	total := int64(math.Round(amount * 100))

	refundParams := &stripe.RefundParams{
		Charge: stripe.String(token),
		Amount: stripe.Int64(total),
	}

	ref, err := refund.New(refundParams)
	return ref, err
}