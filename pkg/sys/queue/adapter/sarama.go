package queueadapter

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/itimky/spindle/pkg/sys/log"
	"github.com/itimky/spindle/pkg/sys/queue"
)

type ConsumerSarama struct {
	handler handler
}

func NewConsumerSarama(
	handler handler,
) *ConsumerSarama {
	return &ConsumerSarama{
		handler: handler,
	}
}

func (c *ConsumerSarama) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerSarama) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerSarama) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := session.Context()
	msgChan := claim.Messages()

	logger, err := log.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("from context: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case saramaMsg := <-msgChan:
			if saramaMsg == nil {
				return nil
			}

			err = c.handler.Handle(ctx, convertSaramaMessage(saramaMsg))
			if err != nil {
				return fmt.Errorf("handle: %w", err)
			}

			session.MarkMessage(saramaMsg, "")

			logger.Infof("message processed: offset %v", saramaMsg.Offset)
		}
	}
}

type BootstrapConsumerSarama struct {
	consumer     partitionConsumerSarama
	bootstrapper bootstrapper
	handler      handler
}

func NewBootstrapConsumerSarama(
	consumer partitionConsumerSarama,
	bootstrapper bootstrapper,
	handler handler,
) *BootstrapConsumerSarama {
	return &BootstrapConsumerSarama{
		consumer:     consumer,
		bootstrapper: bootstrapper,
		handler:      handler,
	}
}

func (c *BootstrapConsumerSarama) Bootstrap(ctx context.Context) error {
	watermark := c.consumer.HighWaterMarkOffset()

	existingMsgs := make([]queue.Message, 0)

	msgChan := c.consumer.Messages()
	errChan := c.consumer.Errors()

LOOP:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		case msg := <-msgChan:
			if msg == nil {
				break LOOP
			}

			existingMsgs = append(existingMsgs, convertSaramaMessage(msg))

			if msg.Offset == watermark-1 {
				break LOOP
			}
		case err := <-errChan:
			return fmt.Errorf("errors: %w", err)
		}
	}

	err := c.bootstrapper.Bootstrap(ctx, existingMsgs)
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}

	return nil
}

func (c *BootstrapConsumerSarama) Consume(ctx context.Context) error {
	msgChan := c.consumer.Messages()
	errChan := c.consumer.Errors()

LOOP:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		case saramaMsg := <-msgChan:
			if saramaMsg == nil {
				break LOOP
			}

			msg := convertSaramaMessage(saramaMsg)

			err := c.handler.Handle(ctx, msg)
			if err != nil {
				return fmt.Errorf("handle: %w", err)
			}
		case err := <-errChan:
			return fmt.Errorf("errors: %w", err)
		}
	}

	return nil
}
