package kafka

import (
	"context"
	"fmt"

	"im-system/internal/config"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func InitProducer(cfg *config.Config) (*Producer, error) {
	brokers := []string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    "messages",
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) PublishMessage(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(cfg *config.Config) (*Consumer, error) {
	brokers := []string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "messages",
		GroupID: "im-group",
	})

	return &Consumer{reader: reader}, nil
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.reader.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
