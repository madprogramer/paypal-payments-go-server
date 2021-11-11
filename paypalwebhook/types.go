package paypalwebhook

import "github.com/plutov/paypal/v4"

type (
	//PayPal Completed Order

	//Any PayPal Event
	paypalevent paypal.AnyEvent
	//paypalevent map[string]interface{}

	//Resource for Checkout Order Approved
	//github.com/plutov/paypal/v4 recommends custom typing based on paypal.Resource
	TransactionResource struct {
		ID                        string                     `json:"id,omitempty"`
		Status                    string                     `json:"status,omitempty"`
		UpdateTime                string                     `json:"update_time,omitempty"`
		CreateTime                string                     `json:"create_time,omitempty"`
		Intent                    string                     `json:"intent,omitempty"`
		PurchaseUnits             []*paypal.PurchaseUnitRequest     `json:"purchase_units,omitempty"`
		Payer                     *paypal.PayerWithNameAndPhone     `json:"payer,omitempty"`
		Links                     []paypal.Link                     `json:"links,omitempty"`
	}

)