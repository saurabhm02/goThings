package constants

import (
	"os"

	"github.com/twilio/twilio-go"
)

var TWILIO_ACCOUNT_SID string = os.Getenv("TWILIO_ACCOUNT_SID")
var TWILIO_AUTH_TOKEN string = os.Getenv("TWILIO_AUTHTOKEN")
var VERIFY_SERVICE_SID string = os.Getenv("TWILIO_SERVICES_ID")
var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username: TWILIO_ACCOUNT_SID,
	Password: TWILIO_AUTH_TOKEN,
})
