//Package pubsub implements publish subscriber patterns for usage in Golang
//go:generate mockgen -source pubsub.go -destination pubsub_mock.go -package pubsub
package pubsub

import (
	"context"
	"time"
)

var (
	client = &Client{Provider: NoopProvider{}}
)

// Client holds a reference to a Provider
type Client struct {
	ServiceName          string
	Provider             Provider
	SubscriberMiddleware []SubscriberMiddleware
}

// SetClient sets the global pubsub client, useful in tests
func SetClient(cli *Client) {
	client = cli
}

// Provider is generic interface for a pub sub provider
type Provider interface {
	Publish(ctx context.Context, topic string, b []byte) error
	Subscribe(topic, subscriberName string, h MsgHandler, deadline time.Duration, autoAck bool)
}

// Subscriber is a service that listens to events and registers handlers
// for those events
type Subscriber interface {
	// Setup is a required method that allows the subscriber service to add handlers
	// and perform any setup if required, this is usually called by pubsub upon start
	Setup(*Client)
}

// Msg is a lile representation of a pub sub message
type Msg struct {
	ID       string
	Metadata map[string]string
	Data     []byte

	Ack  func()
	Nack func()
}

// Handler is a specific callback used for Subscribe in the format of..
// func(ctx context.Context, obj proto.Message, msg *Msg) error
// for example, you can unmarshal a custom type..
// func(ctx context.Context, accounts accounts.Account, msg *Msg) error
// you can also unmarshal a JSON object by supplying any type of interface{}
// func(ctx context.Context, accounts models.SomeJSONAccount, msg *Msg) error
type Handler interface{}

// MsgHandler is the internal or raw message handler
type MsgHandler func(ctx context.Context, m Msg) error

// SubscriberMiddleware is the function signature for injecting middlware
// into subscribers
type SubscriberMiddleware func(opts HandlerOptions, next MsgHandler) MsgHandler
