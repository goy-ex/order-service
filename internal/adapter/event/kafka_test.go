package event_test

import (
	"testing"

	"github.com/goy-ex/order-service/internal/adapter/event"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

func TestNewKafkaOrderProducer(t *testing.T) {
	t.Parallel()

	type testcase struct {
		brokers []string
		check   func(*testing.T, error)
	}

	testcases := map[string]testcase{
		"valid": {
			brokers: []string{"localhost:8080"},
			check: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
		"no_brokers": {
			brokers: nil,
			check: func(t *testing.T, err error) {
				t.Helper()

				require.Error(t, err)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := event.NewKafkaOrderProducer("topic", tc.brokers...)
			tc.check(t, err)
			t.Log(err)
		})
	}
}

func TestKafkaOrderProducer_Produce(t *testing.T) {
	t.Parallel()

	kafkaContainer, err := kafka.Run(
		t.Context(),
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"),
	)

	require.NoError(t, err)

	t.Cleanup(func() {
		err = testcontainers.TerminateContainer(kafkaContainer)
		require.NoError(t, err)
	})

	brokers, err := kafkaContainer.Brokers(t.Context())
	require.NoError(t, err)

	admin, err := kgo.NewClient(kgo.SeedBrokers(brokers...))
	require.NoError(t, err)

	defer admin.Close()

	_, err = kadm.NewClient(admin).CreateTopic(t.Context(), 1, 1, nil, "topic")
	require.NoError(t, err)

	producer, err := event.NewKafkaOrderProducer("topic", brokers...)
	require.NoError(t, err)

	type testcase struct {
		producer *event.KafkaOrderProducer
		events   []*event.Event
		check    func(t *testing.T, err error)
	}

	testcases := map[string]testcase{
		"valid": {
			producer: producer,
			events: []*event.Event{
				{},
			},
			check: func(t *testing.T, err error) {
				t.Helper()

				require.NoError(t, err)
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.producer.Produce(t.Context(), tc.events)
			tc.check(t, err)
		})
	}
}
