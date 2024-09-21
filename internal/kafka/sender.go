package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"

	domain "gitlab.ozon.dev/antonkraeww/ozon-route256-homeworks/internal/models/domain/event"
)

type Sender struct {
	producer sarama.SyncProducer
	topic    string
}

func NewSender(brokers []string, topic string) (*Sender, error) {
	producer, err := NewSyncProducer(brokers)
	if err != nil {
		return nil, err
	}

	return &Sender{producer: producer, topic: topic}, nil
}

func (s *Sender) SendMessage(message domain.Event) error {
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		return err
	}

	_, _, err = s.producer.SendMessage(kafkaMsg)

	return err
}

func (s *Sender) buildMessage(message domain.Event) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
		Key:       sarama.StringEncoder(fmt.Sprint(message.ID)),
	}, nil
}
