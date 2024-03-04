package queueadapter

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/itimky/spindle/pkg/sys/queue"
)

type bootstrapper interface {
	Bootstrap(ctx context.Context, msgs []queue.Message) error
}

type handler interface {
	Handle(ctx context.Context, msg queue.Message) error
}

type partitionConsumerSarama interface {
	Messages() <-chan *sarama.ConsumerMessage
	Errors() <-chan *sarama.ConsumerError
	HighWaterMarkOffset() int64
}
