package services

import (
	"git.nextgencode.io/huyen.vu/freeze-app-rest/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
	"math"
	"os"
	"strconv"
)

func StripeConnectDestinationCharge(token string, accId string, desc string, amount float64) (*stripe.Charge, error) {
	platformFeePercent, _ := strconv.ParseFloat(os.Getenv("PLATFORM_FEE_PERCENTAGE"), 64)
	total := int64(amount*100)
	desAmount := int64(math.Round(amount * (100 - platformFeePercent)))
	params := &stripe.ChargeParams{
		Amount: stripe.Int64(total),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(desc),
		Destination: &stripe.DestinationParams{
			Amount: stripe.Int64(desAmount),
			Account: stripe.String(accId),
		},
	}

	err := params.SetSource(token)

	if err != nil {
		panic(err)
		return nil , err
	}

	return charge.New(params)
}

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

func StripeApplicationFee(token string, accID string, amount float64) (interface{}, error) {
	total := int64(math.Round(amount * 100))
	applicationFee, _ := strconv.ParseInt(os.Getenv("APPLICATION_FEE"), 0, 64	)

	params := &stripe.ChargeParams{
		Amount: stripe.Int64(total),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		ApplicationFee: stripe.Int64(applicationFee*100),
	}

	params.SetStripeAccount(accID)
	_ = params.SetSource(token)

	return charge.New(params)
}

func StripeGetAccountById(id string) (*stripe.Account, error) {
	return account.GetByID(id, nil)
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