package fcm

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FCM interface {
	SendNotification(ctx context.Context, token string, notification NotificationData) error
	SendMulticastNotification(ctx context.Context, tokens []string, notification NotificationData) error
}

type fCMClient struct {
	client *messaging.Client
}

type NotificationData struct {
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Data     map[string]string `json:"data,omitempty"`
	ImageURL string            `json:"imageUrl,omitempty"`
}

func New(serviceJSON string) (FCM, error) {
	opt := option.WithCredentialsFile(serviceJSON)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &fCMClient{client: client}, nil
}

func (f *fCMClient) SendNotification(ctx context.Context, token string, notification NotificationData) error {
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.ImageURL,
		},
		Data: notification.Data,
	}

	_, err := f.client.Send(ctx, message)
	return err
}

func (f *fCMClient) SendMulticastNotification(ctx context.Context, tokens []string, notification NotificationData) error {
	message := &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.ImageURL,
		},
		Data: notification.Data,
	}

	_, err := f.client.SendEachForMulticast(ctx, message)
	return err
}
