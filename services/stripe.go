package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
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

func StripeCreateAccount(merchant models.Merchant) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Country: stripe.String("US"),
		Type: stripe.String(string(stripe.AccountTypeCustom)),
	}

	acc, err := account.New(params)

	if err != nil {
		panic(err)
		return nil, err
	}

	return acc, err
}

func StripeUpdateAccount(merchant models.Merchant) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		SupportPhone: stripe.String(merchant.PhoneNumber),
		Email: stripe.String(merchant.Email),
		SupportEmail: stripe.String(merchant.Email),
	}

	acc, err := account.Update(merchant.StripeID, params)

	if err != nil {
		panic(err)
		return nil, err
	}

	return acc, err
}