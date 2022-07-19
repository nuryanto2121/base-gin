package fcmgetway

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type SendFCM struct {
	Title          string   `json:"title"`
	Body           string   `json:"body"`
	DeviceToken    []string `json:"device_token"`
	PhoneNo        string   `json:"phone_no"`
	Avatar         string   `json:"avatar"`
	Type           string   `json:"type"`
	ConversationID string   `json:"conversation_id"`
}

func (s *SendFCM) SendPushNotification(batchResponse chan *messaging.BatchResponse) error {
	// [START send_multicast]
	// Create a list containing up to 100 registration tokens.
	// This registration tokens come from the client FCM SDKs.

	opt := option.WithCredentialsFile("firebase-sdk-key.json")
	// app, err := firebase.NewApp(context.Background(), nil, opt)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"title":           s.Title,
			"body":            s.Body,
			"type":            s.Type,
			"phone_no":        s.PhoneNo,
			"conversation_id": s.ConversationID,
		},
		Tokens: s.DeviceToken,
		Notification: &messaging.Notification{
			Title: s.Title,
			Body:  s.Body,
		},
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				Title:    s.Title,
				Body:     s.Body,
				Icon:     s.Avatar,
				ImageURL: s.Avatar,
			},
		},
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	batchResponse <- br
	if br.SuccessCount > 0 {
		fmt.Printf("%d messages were sent successfully | %#v\n", br.SuccessCount, s)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.

	// [END send_multicast]

	return nil
}
