package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/hirpc/hrpc/log"
)

type consumerGroupHandler struct {
	ctx context.Context
	h   Handler
}

func (c consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (c consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (c consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := c.h.Handle(c.ctx, msg.Topic, msg.Key, msg.Value, msg.Partition, msg.Offset); err != nil {
			return err
		}
		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}

func RegisterGroupConsumer(ctx context.Context, h Handler, groupID string, topics ...string) {
	g, err := sarama.NewConsumerGroupFromClient(groupID, k.client)
	if err != nil {
		return
	}
	defer g.Close()

	go func(ctx context.Context, h Handler, g sarama.ConsumerGroup, topics ...string) {
		for {
			if err := g.Consume(ctx, topics, consumerGroupHandler{
				ctx: ctx,
				h:   h,
			}); err != nil {
				log.WithFields(ctx).Error(err)
			}
		}
	}(ctx, h, g, topics...)
}

type Handler interface {
	Handle(ctx context.Context, topic string, key, value []byte, partition int32, offset int64) error
}

func Consume(topic string, partition int, f func(m *Message)) error {
	c, err := sarama.NewConsumerFromClient(k.client)
	if err != nil {
		return err
	}
	defer c.Close()
	p, err := c.ConsumePartition(topic, int32(partition), k.opt.OffsetInitial)
	if err != nil {
		return err
	}
	defer p.Close()

	for {
		msg := <-p.Messages()
		f(&Message{
			Topic:     msg.Topic,
			Value:     msg.Value,
			Offset:    msg.Offset,
			Partition: msg.Partition,
			Timestamp: msg.Timestamp,
			Key:       msg.Key,
		})
	}
}
