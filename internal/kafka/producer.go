package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(kafkaBroker string, kafkaTopic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
		Logger:   kafka.LoggerFunc(logf),
	}

	return &Producer{
		w: writer,
	}
}

func (p *Producer) WriteMesage(ctx context.Context, payload []byte) error {
	err := p.w.WriteMessages(ctx, kafka.Message{
		Value: payload,
	})
	if err != nil {
		return err
	}

	return nil
}

func logf(msg string, a ...any) {
	fmt.Printf(msg, a...)
	fmt.Println()
}
