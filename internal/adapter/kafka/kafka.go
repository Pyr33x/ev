package kafka

import (
	"context"
	"fmt"

	"github.com/pyr33x/ev/pkg/config"
	"github.com/segmentio/kafka-go"
)

type Adapter struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func New(cfg *config.Kafka) (*Adapter, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	}

	readerConfig := kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Topic:   cfg.Topic,
		// Partition: cfg.Partition
		GroupID:  cfg.GroupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	}
	reader := kafka.NewReader(readerConfig)

	return &Adapter{
		writer: writer,
		reader: reader,
	}, nil
}

func (a *Adapter) WriteMessage(ctx context.Context, key, value []byte, headers ...kafka.Header) error {
	return a.writer.WriteMessages(ctx, kafka.Message{
		Key:     key,
		Value:   value,
		Headers: headers,
	})
}

func (a *Adapter) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return a.reader.ReadMessage(ctx)
}

func (a *Adapter) Close() error {
	var writerErr, readerErr error

	if a.writer != nil {
		writerErr = a.writer.Close()
	}

	if a.reader != nil {
		readerErr = a.reader.Close()
	}

	if writerErr != nil {
		return writerErr
	}

	return readerErr
}
