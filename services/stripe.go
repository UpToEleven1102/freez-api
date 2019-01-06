package services

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"math"
)

func StripeCharge(token string, des string, amount float64) (interface{}, error) {
	total := int64(math.Round(amount)) * 100

	chargeParams := &stripe.ChargeParams{
		Amount: stripe.Int64(total),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(des),
	}
	chargeParams.SetSource(token)

	ch, err := charge.New(chargeParams)

	return ch, err
}
