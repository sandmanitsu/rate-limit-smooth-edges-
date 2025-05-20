package kafka

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	metric "rate/internal/metrics"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	msgCount    = 100
	workerCount = 30
)

func StartConsumer(kafkaBroker, kafkaTopic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		Topic:    kafkaTopic,
		GroupID:  "msgs-rate",
		MaxBytes: 10e6,
		// Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	})

	ctx := context.Background()

	log.Println("consumer started!")

	msgBatch := make([]kafka.Message, 0, msgCount)
	go func() {
		for {
			m, err := r.FetchMessage(ctx)
			if err != nil {
				log.Println("error fetch message", err)
				continue
			}

			fmt.Println(string(m.Value))

			if len(msgBatch) != msgCount {
				msgBatch = append(msgBatch, m)
				continue
			}

			processedMsgs := processMessage(msgBatch)

			for _, msg := range processedMsgs {
				if err = r.CommitMessages(ctx, msg); err != nil {
					log.Println("error commit message", err)
				}
			}
		}
	}()
}

func processMessage(msgs []kafka.Message) []kafka.Message {
	inputCh := make(chan kafka.Message)
	outputCh := make(chan kafka.Message)
	wg := &sync.WaitGroup{}

	output := make([]kafka.Message, 0, len(msgs))

	go func() {
		defer close(inputCh)

		for _, m := range msgs {
			inputCh <- m
		}
	}()

	go func() {
		for _ = range workerCount {
			wg.Add(1)

			go messageWorker(wg, inputCh, outputCh)
		}

		wg.Wait()
		close(outputCh)
	}()

	for res := range outputCh {
		output = append(output, res)
	}

	return output
}

func messageWorker(wg *sync.WaitGroup, inCh <-chan kafka.Message, outCh chan<- kafka.Message) {
	start := time.Now()
	defer func() {
		wg.Done()
		metric.ObserveCodeStatus(1, time.Since(start))
	}()

	for inData := range inCh {
		// some work...
		randomMs := rand.Intn(1000) + 100
		time.Sleep(time.Millisecond * time.Duration(randomMs))

		log.Println("work done!")

		outCh <- inData
	}
}
