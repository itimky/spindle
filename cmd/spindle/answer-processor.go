package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	answerprocessor "github.com/itimky/spindle/pkg/domain/answer-processor"
	matchfacade "github.com/itimky/spindle/pkg/facade/match"
	kafkahandler "github.com/itimky/spindle/pkg/handler/kafka"
	answerprocessorstorage "github.com/itimky/spindle/pkg/storage/answer-processor"
	"github.com/itimky/spindle/pkg/sys/log"
	"github.com/itimky/spindle/pkg/sys/queue"
	queueadapter "github.com/itimky/spindle/pkg/sys/queue/adapter"
	"github.com/itimky/spindle/pkg/sys/run"
)

type AnswerProcessor struct {
	matrixInMem          *answerprocessorstorage.MatrixInMem
	answersConsumer      *queueadapter.ConsumerSarama
	logger               log.Logger
	brokers              []string
	answersTopic         string
	group                string
	handler              *queue.Handler
	kafkaConsumerBackoff time.Duration
}

func NewAnswerProcessor(
	cfg Config,
	logger log.Logger,
) *AnswerProcessor {
	matrixInMemStore := answerprocessorstorage.NewMatrixInMem(
		run.NewReadiness(),
		&sync.RWMutex{},
	)
	repo := answerprocessorstorage.NewComposite(
		matrixInMemStore,
		nil,
	)
	answerProcessor := answerprocessor.NewAnswerProcessor(repo)
	matchFacade := matchfacade.NewMatchFacade(answerProcessor)
	kafkaHandler := kafkahandler.NewHandler(matchFacade)
	queueHandler := queue.NewHandler(kafkaHandler, queue.MiddlewareRecover)
	consumer := queueadapter.NewConsumerSarama(queueHandler)

	return &AnswerProcessor{
		matrixInMem:          matrixInMemStore,
		answersConsumer:      consumer,
		logger:               logger,
		answersTopic:         cfg.AnswersTopic,
		brokers:              cfg.KafkaBrokers,
		group:                cfg.AnswerProcessorGroup,
		handler:              queueHandler,
		kafkaConsumerBackoff: cfg.KafkaConsumerBackoff,
	}
}

func (c *AnswerProcessor) RunConsumer(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(
		c.brokers,
		c.group,
		config,
	)
	if err != nil {
		return fmt.Errorf("new consumer group: %w", err)
	}

	defer func() {
		if err := group.Close(); err != nil {
			c.logger.Errorf("group close: %s", err)
		}
	}()

	c.logger.Info("Consumer has been started")

LOOP:
	for {
		err = group.Consume(ctx, []string{c.answersTopic}, c.answersConsumer)
		if err != nil {
			c.logger.Errorf("consume: %s", err)

			time.Sleep(c.kafkaConsumerBackoff)
		}

		select {
		case <-ctx.Done():
			break LOOP
		case err = <-group.Errors():
			c.logger.Errorf("consume: %s", err)

			select {
			case <-ctx.Done():
				break LOOP
			case <-time.After(c.kafkaConsumerBackoff):
			}
		default:
		}
	}

	c.logger.Info("Consumer has been stopped")

	return nil
}
