package push_notification

import (
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/wizzldev/chat/pkg/configs"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

type pushNotification struct {
	init   bool
	client *messaging.Client
}

var PushNotification = &pushNotification{
	init:   false,
	client: nil,
}

func (p *pushNotification) Init() error {
	if p.init {
		return nil
	}

	p.init = true

	decodedKey, err := p.getKey()
	if err != nil {
		return err
	}

	opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	app, err := firebase.NewApp(context.Background(), nil, opts...)
	if err != nil {
		return err
	}

	fcmClient, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	p.client = fcmClient

	return nil
}

func (p *pushNotification) Send(tokens []string, title string, body string, imageURL string) error {
	if !p.init {
		return nil
	}
	_, err := p.client.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imageURL,
		},
		Tokens: tokens,
	})
	return err
}

func (*pushNotification) getKey() ([]byte, error) {
	return base64.StdEncoding.DecodeString(configs.Env.FirebaseAuthKey)
}
