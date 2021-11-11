// Working example of a basic Go server listening on a Paypal Webhook
package main

import (
	"net/http"
    "fmt"
    "log"
    "os"

    "example.com/paypalwebhook"
    "github.com/plutov/paypal/v4"
)

var (
	//Pointer to a Paypal Client (from plutov/paypal on GitHub)
	c *paypal.Client
)

//
func main() {
	log.Print("\033[33mPlease remember to update the Webhook URL \033[35mfrom your Paypal Developer account\033[33m. Make sure that you are not on an ngrok address and that your endpoint is connected to the internet.\033[0m")
	log.Print("\033[33mAlso check your PayPal Account to make sure that your Webhook is the same version as your API Mode. Is it set to \033[35mSandbox\033[33m or \033[35mLive\033[33m? Is it spelt correctly?\033[0m\n")
	
	c, _ = paypalwebhook.GetPayPalClient()
	//c, _ = paypalwebhook.GetPayPalClientWith(clientID,secretID,apiMode,webhookID)
	
	if(enablePayPalClientOutput){c.SetLog(os.Stdout)}
	//http.HandleFunc("/oldOrder", oldOrder)
	http.HandleFunc(webhookpath, webhookHandler)
	http.HandleFunc("/", defaultHandler)
	log.Fatal(http.ListenAndServe(defaultport, nil))
}

//Default Handler
//Place Holder.
//
//TODO: Add some info text/link to a tutorial here
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("HI!");
	fmt.Fprint(w, "HI!")
}

//Webhook Handler
//Passes Paypal Client Object to the PaypalWebhook function
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	//log.Print("Webhook?")
	paypalwebhook.PaypalWebhook(c,w,r)
}


/*func oldOrder(w http.ResponseWriter, r *http.Request) {
}*/
