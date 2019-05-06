package services

import (
	"fmt"
	"git.nextgencode.io/huyen.vu/freez-app-rest/models"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/account"
	"github.com/stripe/stripe-go/balance"
	"github.com/stripe/stripe-go/card"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/refund"
	"math"
	"os"
	"strconv"
	"time"
)

func StripeCharge(token string, des string, amount float64) (*stripe.Charge, error) {
	total := int64(math.Round(amount * 100))

	chargeParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(total),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
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

func StripeGetAccountBalance(accId string) (*stripe.Balance, error) {
	params := &stripe.BalanceParams{
		Params: stripe.Params{StripeAccount: stripe.String(accId)},
	}

	return balance.Get(params)
}

//stripe connect

func StripeConnectDestinationCharge(token string, accId string, desc string, amount float64) (*stripe.Charge, error) {
	platformFeePercent, _ := strconv.ParseFloat(os.Getenv("PLATFORM_FEE_PERCENTAGE"), 64)
	total := int64(amount * 100)
	desAmount := int64(math.Round(amount * (100 - platformFeePercent)))
	params := &stripe.ChargeParams{
		Amount:      stripe.Int64(total),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String(desc),
		TransferData: &stripe.ChargeTransferDataParams{
			Amount: stripe.Int64(desAmount),
			Destination: stripe.String(accId),
		},
	}

	err := params.SetSource(token)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return charge.New(params)
}

func StripeConnectCreateAccount(merchant models.Merchant, ipAdd string) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Country:               stripe.String("US"),
		Email:                 stripe.String(merchant.Email),
		BusinessType: stripe.String("individual"),
		RequestedCapabilities: []*string{stripe.String("platform_payments")},
		DefaultCurrency:       stripe.String(string(stripe.CurrencyUSD)),
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Token:    stripe.String(merchant.StripeID),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
		},
		Type: stripe.String(string(stripe.AccountTypeCustom)),
		TOSAcceptance: &stripe.AccountTOSAcceptanceParams{
			Date: stripe.Int64(time.Now().Unix()),
			IP:   stripe.String(ipAdd),
		},
	}

	acc, err := account.New(params)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return acc, err
}


func StripeConnectCreateAccountWithEntityVerification(merchant models.Merchant, ipAdd string, entityInfo models.MerchantEntityInformation) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		Country:               stripe.String("US"),
		Email:                 stripe.String(merchant.Email),
		BusinessType: stripe.String("individual"),
		RequestedCapabilities: []*string{stripe.String("platform_payments")},
		DefaultCurrency:       stripe.String(string(stripe.CurrencyUSD)),
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Token:    stripe.String(merchant.StripeID),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
		},
		Type: stripe.String(string(stripe.AccountTypeCustom)),
		TOSAcceptance: &stripe.AccountTOSAcceptanceParams{
			Date: stripe.Int64(time.Now().Unix()),
			IP:   stripe.String(ipAdd),
		},
		Individual: &stripe.PersonParams{
			FirstName: stripe.String(entityInfo.EntityInformation.FirstName),
			LastName: stripe.String(entityInfo.EntityInformation.LastName),
			Address: &stripe.AccountAddressParams{
				Line1: stripe.String(entityInfo.EntityInformation.Line1),
				City: stripe.String(entityInfo.EntityInformation.City),
				State: stripe.String(entityInfo.EntityInformation.State),
			},
			DOB: &stripe.DOBParams{
				Day: stripe.Int64(entityInfo.EntityInformation.DobDay),
				Month: stripe.Int64(entityInfo.EntityInformation.DobMonth),
				Year: stripe.Int64(entityInfo.EntityInformation.DobYear),
			},
			SSNLast4: stripe.String(entityInfo.EntityInformation.Last4Ssn),
		},
		BusinessProfile: &stripe.AccountBusinessProfileParams{
			URL: stripe.String(entityInfo.EntityInformation.BusinessWebsite),
		},
	}

	acc, err := account.New(params)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return acc, err
}

func StripeConnectGetAccountById(id string) (*stripe.Account, error) {
	return account.GetByID(id, nil)
}

func StripeConnectAddDebitCard(accId string, token string) (*stripe.Account, error) {
	params := &stripe.AccountParams{
		ExternalAccount: &stripe.AccountExternalAccountParams{
			Token: stripe.String(token),
		},
	}

	return account.Update(accId, params)
}

func StripeConnectCreateDebitCard(accId string, token string) (*stripe.Card, error) {
	params := &stripe.CardParams{
		Account: stripe.String(accId),
		Token:   stripe.String(token),
	}

	return card.New(params)
}

func StripeConnectDeleteDebitCard(accId string, cardId string) (*stripe.Card, error) {
	params := &stripe.CardParams{
		Account: stripe.String(accId),
	}

	return card.Del(cardId, params)
}

func StripeConnectGetCardListByStripeId(stripeId string) (cards []*stripe.Card, err error) {
	params := &stripe.CardListParams{
		Account: stripe.String(stripeId),
	}

	params.Filters.AddFilter("limit", "", "3")
	i := card.List(params)
	for i.Next() {
		c := i.Card()
		cards = append(cards, c)
	}

	return cards, err
}

func StripeConnectMakeDefaultCurrencyDebitCard(stripeId string, cardId string) (*stripe.Card, error) {
	params := &stripe.CardParams{
		Account:            stripe.String(stripeId),
		DefaultForCurrency: stripe.Bool(true),
	}

	return card.Update(cardId, params)
}
