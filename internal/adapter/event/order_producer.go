package event

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaOrderProducer struct {
	topic  string
	client *kgo.Client
}

func NewKafkaOrderProducer(topic string, brokers ...string) (*KafkaOrderProducer, error) {
	client, err := kgo.NewClient(kgo.SeedBrokers(brokers...))
	if err != nil {
		return nil, fmt.Errorf("failed to create kgo client: %w", err)
	}

	return &KafkaOrderProducer{
		topic:  topic,
		client: client,
	}, nil
}

func (op *KafkaOrderProducer) Produce(ctx context.Context, events []*Event) error {
	records := make([]*kgo.Record, 0, len(events))

	for _, e := range events {
		records = append(
			records,
			&kgo.Record{
				Topic: op.topic,
				Key:   []byte(e.PartitionKey),
				Value: e.Payload,
				Headers: []kgo.RecordHeader{
					{
						Key:   "event_type",
						Value: []byte(e.EventType.String()),
					},
				},
			},
		)
	}

	results := op.client.ProduceSync(ctx, records...)

	err := results.FirstErr()
	if err != nil {
		return fmt.Errorf("failed to produce record: %w", err)
	}

	return nil
}
