package main

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	systemfacade "github.com/itimky/spindle/pkg/facade/system"
	kafkahandler "github.com/itimky/spindle/pkg/handler/kafka"
	answerprocessorstorage "github.com/itimky/spindle/pkg/storage/answer-processor"
	"github.com/itimky/spindle/pkg/sys/log"
	queueadapter "github.com/itimky/spindle/pkg/sys/queue/adapter"
)

type System struct {
	brokers             []string
	answerMatricesTopic string
	logger              log.Logger
	facade              *systemfacade.SystemFacade
	handler             *kafkahandler.WeightMatrixBootstrapHandler
}

func NewSystem(
	cfg Config,
	logger log.Logger,
	matrixInMem *answerprocessorstorage.MatrixInMem,
) *System {
	facade := systemfacade.NewSystemFacade(matrixInMem)
	handler := kafkahandler.NewWeightMatrixBootstrapHandler(facade)

	return &System{
		brokers:             cfg.KafkaBrokers,
		answerMatricesTopic: cfg.WeightMatricesTopic,
		logger:              logger,
		facade:              facade,
		handler:             handler,
	}
}

func (c *System) RunBootstrapConsumer(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	var partition int32 // default partition 0

	consumer, err := sarama.NewConsumer(
		c.brokers,
		config,
	)
	if err != nil {
		return fmt.Errorf("new consumer: %w", err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			c.logger.Errorf("consumer close: %s", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(c.answerMatricesTopic, partition, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("consume partition: %w", err)
	}

	defer func() {
		if partitionConsumer.Close() != nil {
			c.logger.Errorf("partition consumer close: %s", err)
		}
	}()

	bootstrapConsumer := queueadapter.NewBootstrapConsumerSarama(
		partitionConsumer,
		c.handler,
		c.handler,
	)

	err = bootstrapConsumer.Bootstrap(ctx)
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}

	c.logger.Info("Consumer has been started")

	err = bootstrapConsumer.Consume(ctx)
	if err != nil {
		return fmt.Errorf("consumer adapter run: %w", err)
	}

	return nil
}
