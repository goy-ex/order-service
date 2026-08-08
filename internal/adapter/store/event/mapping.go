package event

import (
	"github.com/goy-ex/order-service/internal/adapter/event"
	pgxsqlc "github.com/goy-ex/order-service/internal/adapter/store/sqlc/pgx"
	"github.com/goy-ex/order-service/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func mapSelectOrderEventBatchRows(rows []pgxsqlc.SelectOrderEventBatchRow) ([]*event.Event, error) {
	events := make([]*event.Event, 0, len(rows))

	for _, r := range rows {
		events = append(events, &event.Event{
			ID:            r.ID.Bytes,
			AggregateID:   r.AggregateID.Bytes,
			AggregateType: string(r.AggregateType),
			PartitionKey:  r.PartitionKey.String,
			EventType:     domain.EventTypeFromString(string(r.EventType)),
			Payload:       r.Payload,
		})
	}

	return events, nil
}

func mapBatchToUUIDs(batch []*event.Event) []pgtype.UUID {
	uuids := make([]pgtype.UUID, len(batch))

	for _, e := range batch {
		uuids = append(uuids, pgtype.UUID{Bytes: e.ID, Valid: true})
	}

	return uuids
}
