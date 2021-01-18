// Package v2 provides a client with included tracing capabilities.
package v2

import (
	"context"
	"errors"
	"fmt"

	"github.com/beatlabs/patron/correlation"
	patronerrors "github.com/beatlabs/patron/errors"
	"github.com/beatlabs/patron/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/streadway/amqp"
)

const (
	publisherComponent = "amqp-publisher"
)

// Publisher defines a RabbitMQ publisher with tracing instrumentation.
type Publisher struct {
	cfg        *amqp.Config
	connection *amqp.Connection
	channel    *amqp.Channel
}

// New constructor.
func New(url string, oo ...OptionFunc) (*Publisher, error) {
	if url == "" {
		return nil, errors.New("url is required")
	}

	var err error
	pub := &Publisher{}

	for _, option := range oo {
		err = option(pub)
		if err != nil {
			return nil, err
		}
	}

	var conn *amqp.Connection

	if pub.cfg == nil {
		conn, err = amqp.Dial(url)
	} else {
		conn, err = amqp.DialConfig(url, *pub.cfg)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, patronerrors.Aggregate(fmt.Errorf("failed to open channel: %w", err), conn.Close())
	}

	pub.connection = conn
	pub.channel = ch
	return pub, nil
}

// Publish a message to a exchange.
func (tc *Publisher) Publish(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	sp, _ := trace.ChildSpan(ctx, trace.ComponentOpName(publisherComponent, exchange),
		publisherComponent, ext.SpanKindProducer, opentracing.Tag{Key: "exchange", Value: exchange})

	if msg.Headers == nil {
		msg.Headers = amqp.Table{}
	}

	c := amqpHeadersCarrier(msg.Headers)
	err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, c)
	if err != nil {
		return fmt.Errorf("failed to inject tracing headers: %w", err)
	}
	msg.Headers[correlation.HeaderID] = correlation.IDFromContext(ctx)

	err = tc.channel.Publish(exchange, key, mandatory, immediate, msg)
	trace.SpanComplete(sp, err)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close the channel and connection.
func (tc *Publisher) Close() error {
	return patronerrors.Aggregate(tc.channel.Close(), tc.connection.Close())
}

type amqpHeadersCarrier map[string]interface{}

// Set implements Set() of opentracing.TextMapWriter.
func (c amqpHeadersCarrier) Set(key, val string) {
	c[key] = val
}
