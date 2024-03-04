package kafkacontract

import "github.com/itimky/spindle/pkg/sys/queue"

const (
	MessageTypeAnswerV1       queue.MessageType = "answer/v1"
	MessageTypeWeightMatrixV1 queue.MessageType = "weight-matrix/v1"
)

type AnswerV1 struct {
	PersonID   string `json:"personId"`
	QuestionID string `json:"questionId"`
	AnswerID   string `json:"answerId"`
}

type WeightMatrixV1 struct {
	QuestionID string                       `json:"questionId"`
	Matrix     map[string]map[string]string `json:"matrix"`
}
