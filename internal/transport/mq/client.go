package mq_client

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
	"github.com/zardan4/petition-audit-rabbitmq/pkg/core/audit"
)

type Client struct {
	ch   *amqp.Channel
	logQ amqp.Queue
}

func NewClient(mqUsername, mqPassword string) (*Client, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@localhost:15672/", mqUsername, mqPassword))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// declare logs queue
	logsQ, err := ch.QueueDeclare(
		viper.GetString("queues.logs"),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		ch:   ch,
		logQ: logsQ,
	}, nil
}

func (c *Client) Close() error {
	return c.ch.Close()
}

func (c *Client) SendLogRequest(ctx context.Context, req audit.LogItem) error {
	sendingBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = c.ch.PublishWithContext(ctx,
		"",
		c.logQ.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        sendingBody,
		},
	)

	return err
}
