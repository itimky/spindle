package queueadapter

import (
	"github.com/IBM/sarama"
	"github.com/itimky/spindle/pkg/sys/queue"
)

func convertSaramaMessage(msg *sarama.ConsumerMessage) queue.Message {
	var msgType queue.MessageType

	for _, header := range msg.Headers {
		if string(header.Key) == "type" {
			if header.Value != nil {
				msgType = queue.MessageType(header.Value)
			}

			break
		}
	}

	return queue.Message{
		Type: msgType,
		Data: msg.Value,
	}
}
