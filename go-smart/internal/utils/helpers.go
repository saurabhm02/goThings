package utils

import (
	"errors"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	twilioClient     *twilio.RestClient
	verifyServiceSID string
)

func init() {
	// Load Twilio credentials
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	verifyServiceSID = os.Getenv("VERIFY_SERVICE_SID")

	if verifyServiceSID == "" {
		log.Fatal("Twilio VERIFY_SERVICE_SID is not set")
	}
}

func GenerateOtp() int {
	rand.Seed(time.Now().UnixNano())
	return 100000 + rand.Intn(900000)
}

func SendMail(to, subject, msg string) error {
	from := os.Getenv("MAIL")
	password := os.Getenv("PASSWD")
	host := "smtp.gmail.com"
	port := "587"

	body := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" + msg + "\r\n")

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, body)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Successfully sent email to ****%s", to[len(to)-4:])
	return nil
}

func SendSms(phone string) (string, error) {
	if verifyServiceSID == "" {
		return "", errors.New("Verify Service SID is not set")
	}

	to := phone
	if !strings.HasPrefix(phone, "+") {
		to = "+91" + phone
	}

	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := twilioClient.VerifyV2.CreateVerification(verifyServiceSID, params)
	if err != nil {
		log.Printf("Twilio SMS failed: %v", err)
		return "", errors.New("OTP failed to send")
	}

	log.Printf("Sent verification SID for phone ****%s", to[len(to)-4:])
	return *resp.Sid, nil
}

func CheckOtp(phone, code string) error {
	to := phone
	if !strings.HasPrefix(phone, "+") {
		to = "+91" + phone
	}

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := twilioClient.VerifyV2.CreateVerificationCheck(verifyServiceSID, params)
	if err != nil {
		log.Printf("OTP verification failed: %v", err)
		return errors.New("Invalid OTP or expired")
	}

	if *resp.Status == "approved" {
		log.Printf("OTP approved for phone ****%s", phone[len(phone)-4:])
		return nil
	}

	return errors.New("Invalid OTP")
}
