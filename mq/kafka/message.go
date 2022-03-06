package kafka

import "time"

// Message represents the kafka message
type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Offset    int64
	Partition int32
	Timestamp time.Time
}

func NewProduceMessage(value []byte) *Message {
	return &Message{
		Value: value,
	}
}
