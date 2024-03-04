package main

import (
	"time"
)

type Config struct {
	KafkaBrokers         []string      `env:"KAFKA_BROKERS"          envDefault:"kafka:9092"`
	KafkaConsumerBackoff time.Duration `env:"KAFKA_CONSUMER_BACKOFF" envDefault:"5s"`

	WeightMatricesTopic  string `env:"WEIGHT_MATRICES_TOPIC"  envDefault:"yarn.weight-matrices"`
	AnswersTopic         string `env:"ANSWERS_TOPIC"          envDefault:"yarn.answers"`
	AnswerProcessorGroup string `env:"ANSWER_PROCESSOR_GROUP" envDefault:"spindle.answer-processor"`
}
