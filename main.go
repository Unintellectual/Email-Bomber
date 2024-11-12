package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"mime"
	"net/mail"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var gmailScopes = []string{gmail.GmailSendScope}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func sendEmail(srv *gmail.Service, to string, subject string, body string) error {
	var message gmail.Message

	// Create the recipient and sender email addresses
	emailTo := mail.Address{Name: "", Address: to} // You can leave Name as "" if you don't need to include the name
	emailFrom := mail.Address{Name: "Kalbo Kobu", Address: "comjoed00509@gmail.com"}

	// Create the email headers
	header := make(map[string]string)
	header["From"] = emailFrom.String()
	header["To"] = emailTo.String()
	header["Subject"] = mime.QEncoding.Encode("utf-8", subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	// Build the email body
	var msg strings.Builder
	for k, v := range header {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body)

	// Encode and send the email
	raw := base64.URLEncoding.EncodeToString([]byte(msg.String()))
	message.Raw = raw

	_, err := srv.Users.Messages.Send("me", &message).Do()
	return err
}

func main() {
	ctx := context.Background()

	// Load credentials from the "credentials.json" file in the root folder
	credentialsFile := "credentials.json"
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Parse the credentials JSON
	config, err := google.ConfigFromJSON(b, gmailScopes...)
	if err != nil {
		log.Fatalf("Unable to parse client secret JSON: %v", err)
	}
	client := getClient(config)

	// Create Gmail API service
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// Define email details
	to := "comdamnsdunnns@gmail.com"
	subject := "Test Email"
	body := "Kobu Kalbo kalbo"

	// Send email (repeat as needed)
	for i := 0; i < 100; i++ {
		if err := sendEmail(srv, to, subject, body); err != nil {
			log.Printf("Error sending email #%d: %v\n", i+1, err)
		} else {
			fmt.Printf("Email #%d sent successfully!\n", i+1)
		}
	}
}
