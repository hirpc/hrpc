package kafka

import "github.com/Shopify/sarama"

func ProduceAsync(topic string, msg Message) error {
	p, err := sarama.NewAsyncProducerFromClient(k.client)
	if err != nil {
		return err
	}
	defer p.Close()

	p.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(msg.Key),
		Value: sarama.StringEncoder(msg.Value),
	}
	return nil
}

func Produce(topic string, msgs ...Message) error {
	p, err := sarama.NewSyncProducerFromClient(k.client)
	if err != nil {
		return err
	}
	defer p.Close()

	var sMessage []*sarama.ProducerMessage
	for _, msg := range msgs {
		sMessage = append(sMessage, &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(msg.Key),
			Value: sarama.ByteEncoder(msg.Value),
		})
	}
	return p.SendMessages(sMessage)
}
