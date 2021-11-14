package main

//Payment Settings
const(
	//Client ID should match YOUR APP'S CLIENT ID
	//sb is a placeholder for testing (operated by paypal)
	clientid = "sb"
	//min amount to pay
	minamount = 30
	//max amount to pay
	//ignored if <= 0
	maxamount = 0
)

//Server Settings
const(
	//Default Server Port
	defaultport = ":8080"
	//Webhook Path
	webhookpath = "/webhook/paypal"
	//Enable/Disable PayPal Client Output to Stdout
	enablePayPalClientOutput = false
)