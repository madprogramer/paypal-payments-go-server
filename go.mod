module github.com/madprogramer/paypal-payments-go-server

go 1.17

require (
	example.com/paypalwebhook v0.0.0-00010101000000-000000000000
	github.com/plutov/paypal/v4 v4.4.1
)

replace example.com/paypalwebhook => ../paypal-payments-go-server/paypalwebhook
