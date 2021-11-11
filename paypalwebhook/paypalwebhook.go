// paypalwebhook provides functions for listening to and verifying PayPal Webhooks.
package paypalwebhook

import (
    "net/http"
    "time"
    "io/ioutil"
    "log"
    "encoding/json"
    "strconv"
    "bytes"
    "context"
    "errors"
    "github.com/plutov/paypal/v4"
)

var (
	//Timestamp of when to renew Access Token
	renewAccessWhen time.Time

	//Default webhook
	orderApprovedWebhookID string
)

const (
	//Error message to check secrets.go and secrets_go_here.go
	//TODO: Link to tutorial?
	packageSecretsError = "Have you remembered to set up your \033[35m`secrets.go`\033[0m file? See \033[35m`secrets_go_here.go`\033[0m for instructions.\nAlternatively, switch your GetPayPalClient function for GetPayPalClientWith which allows you to pass arguments directly."
)

/*
Returns a Paypal Client for making requests to the PayPal API, and a timestamp for when to renew its access token.
ClientID, Secret ID etc. are all passed from paypalwebhook/secrets.go.

To pass these directly, call `GetPayPalClientWith` instead.
*/
func GetPayPalClient() (*paypal.Client, time.Time) {
	if clientID==""{
		panic(errors.New("\033[31mError: PayPal Client ID is undefined.\n\033[0m"+packageSecretsError))
	} else if secretID==""{
		panic(errors.New("\033[31mError: PayPal Secret is undefined.\n\033[0m"+packageSecretsError))
	} else if apiMode==""{
		panic(errors.New("\033[31mError: PayPal API Mode is undefined.\n\033[0m"+packageSecretsError))
	} else if webhookID==""{
		panic(errors.New("\033[31mError: PayPal Webhook ID is undefined.\n\033[0m"+packageSecretsError))
	}
	return GetPayPalClientWith(clientID,secretID,apiMode,webhookID)
}

/*
Returns a Paypal Client for making requests to the PayPal API, and a timestamp for when to renew its access token.
Takes as input four strings: `clientID`, `secretID`, `apiMode` and `webhookID`.

To pass these directly, call `GetPayPalClientWith` instead.
*/
func GetPayPalClientWith(clientID string,secretID string, apiMode string ,webhookID string) (*paypal.Client, time.Time) {

	log.Print("\033[32mGenerating new Access Token for PayPal REST API...\033[0m")
	log.Print("\033[33mPlease remember to update \033[35mclient/secretID in secrets.go\033[33m \033[1;33mand to swap\033[0m\033[33m \033[35mpaypal.APIBaseSandBox\033[33m with \033[35mpaypal.APIBaseLive\033[33m before going live.\033[0m\n")
	//var t *paypal.TokenResponse;

	//Get a new PayPal client
    c, err := paypal.NewClient(clientID, secretID, apiMode)
    if(err != nil) {panic(err)}

	//Update access token in the future using checkAccessTokenRenewalTime()
	var t *paypal.TokenResponse
	t, err = c.GetAccessToken(context.Background())
	if(err != nil) {panic(err)}
	//log.Print(t)
	//log.Print(t.Token)
	renewAccessWhen = time.Now().Add(time.Duration(t.ExpiresIn) * time.Second)

	//Set Webhook ID
	//Change this line if you want to enable multiple webhookIDs
	orderApprovedWebhookID = webhookID

	//The Paypal client object c is returned.
	//If you are only going to use c for webhooks, the paypalwebhook.go built-ins will handle token expiration etc. automatically
	//If not, make sure to call the GetAccessTokenRenewalTime() function to check when to renew access time
	//Timestamp for when to renew access token is also returned

	return c, renewAccessWhen;
}

/*
Get Access Token Renewal Time
*/
func GetAccessTokenRenewalTime()(time.Time){
	return renewAccessWhen;
}

/*
Check Time and Renew Access Token if it's about to expire!
Returns an error if it has occured.
*/
func checkAccessTokenRenewalTime(c *paypal.Client)(err error){
	//Step 1: Check current time

	currentTime := time.Now()
	//Step 2: Check c.tokenExpiresAt
	//Step 3: Compare diff

	diff := renewAccessWhen.Sub(currentTime)
	//Step 4: If too close or past, _, err = c.GetAccessToken(context.Background())

	if (diff.Seconds() < 300){
		log.Print("\033[33mAccess token to PayPal is about to expire\n\033[32mAttempting to generate new accessToken...\033[0m")
		_, err = c.GetAccessToken(context.Background())
	}
	//Step 5: Check for errors

	return err
}

//Force renew acces token
//Returns new token's expiration time and an error if it has occured.
func RenewAccessToken(c *paypal.Client)(time.Time, error){
	var err error;
	log.Print("\033[32mAttempting to generate new accessToken...\033[0m")
	_, err = c.GetAccessToken(context.Background())
	//Check for errors
	return renewAccessWhen, err
}

/*
Return the account and amount for a transaction resource. An error is returned if it occurs.
Account is passed through purchase unit description, to account* for cases where a paypal account for someone else is used to make a purchase.

* (pun intended) 
*/
func getAccountAndAmount(r TransactionResource)(string, float64, error){
	//Unparse Order Approved Object
	account := r.PurchaseUnits[0].Description
	amountstr := r.PurchaseUnits[0].Amount.Value;
	if amount, err := strconv.ParseFloat(amountstr, 64); err == nil{
		return account, amount, nil
	}else {
		return "", 0, err;
	}
}

/*
Verify Paypal Events (namely CHECKOUT.ORDER.APPROVEDs)
Returns boolean indicating success or failure and, possibly, an error
*/
func verifyPaypalEvent(c *paypal.Client, r *http.Request, whID string, event_id string ) (bool, error) {
	
	//Check For New Access Token First
	err := checkAccessTokenRenewalTime(c)
	if(err != nil) {return false, err}

	log.Print("\033[33mVerifying Webhook Event: ", event_id, "...\033[0m")
	status, err_verifywebhook := c.VerifyWebhookSignature(context.Background(),r,webhookID)
	if(err_verifywebhook != nil) {return false, err_verifywebhook}

	//CHECK IF STATUS IS SUCCESS OR FAILUREx
	return status.VerificationStatus=="SUCCESS", nil
}

// PaypalWebhook is a handler for picking up Webhook transaction updates from PayPal
//
func PaypalWebhook(c *paypal.Client, w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	defer r.Body.Close()

	if r.Method == "POST" {

		// Need to pass request onto PayPal for verification
		// Read the content for local processing, and then reset the io.ReadCloser
		var jsonBytes []byte
		if r.Body != nil {
		  jsonBytes, _ = ioutil.ReadAll(r.Body)
		}
		// Restore the io.ReadCloser to its original state
		r.Body = ioutil.NopCloser(bytes.NewBuffer(jsonBytes))

		// Decode JSON
		var transaction paypalevent

		d := json.NewDecoder(bytes.NewReader(jsonBytes))
		d.DisallowUnknownFields() // catch unwanted fields
		err := d.Decode(&transaction)
		
		if err != nil {
			//panic(err)
			log.Print(err)
			log.Print("\033[31m Received POST, unable to decode into JSON. See log.\033[0m");
			return
		}
		
		//Attempt to check event type of transaction
		event_type := transaction.Event.EventType
		if event_type == "" {log.Print("\033[31mReceived Malformed POST request. Ignoring.\033[0m"); return}

		event_id := transaction.Event.ID
		if event_id == "" {log.Print("\033[31mCannot find event id for PayPal Event.\033[0m"); return}

		log.Print("\033[34mNew ",event_type," Transaction Received: ", event_id ,"\033[0m")

		//We expect to receive a CHECKOUT ORDER APPROVED EVENT
		if event_type == "CHECKOUT.ORDER.APPROVED" {
		  
		  	//Get Account and Amount
		  	var tr TransactionResource;
			err := json.Unmarshal(transaction.Resource, &tr)
			if err != nil {
				log.Print(err)
				log.Print("\033[31mUnable to decode resource into JSON. See log.\033[0m");
				return;
			}

			var account string;
			var amount float64;
			account, amount, err = getAccountAndAmount(tr)

			if err != nil{
		    	log.Print(err)
		    	log.Print("\033[31mAn error occured while reading Webhook Content for ",event_id,", see log.\033[0m");
		    	return;
		    }

		    //Verify
		    var v bool;

		    v, err = verifyPaypalEvent(c,r,orderApprovedWebhookID,event_id)
		    if v {
		    	//DB Update
		    	log.Print("\033[35mDB Update Required\n","Add ", amount, " to ", account,".\033[0m");

		    } else{
		    	log.Print(err)
		    	log.Print("\033[31mVerification error received for ",event_id,", see log.\033[0m");
		    	return;
		    }

		} else { 
			//If you receive this message, it means you need to implement a new branching condition for 
			//}else if event_type == "XXX"{}
			log.Print("Unrecognized event type: ", event_type)
		}
	}
}