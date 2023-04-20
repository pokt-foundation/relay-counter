package postgresdriver

import (
	"context"
	"time"

	"github.com/pokt-foundation/relay-counter/types"
)

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, &time.Location{})
}

func (d *PostgresDriver) WriteRelayCount(ctx context.Context, count types.RelayCount) error {
	now := time.Now()

	return d.InsertRelayCount(ctx, InsertRelayCountParams{
		AppPublicKey: count.AppPublicKey,
		Day:          truncateToDay(count.Day),
		Success:      int32(count.Success),
		Error:        int32(count.Error),
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func (d *PostgresDriver) ReadRelayCounts(ctx context.Context, from, to time.Time) ([]types.RelayCount, error) {
	dbCounts, err := d.SelectRelayCounts(ctx, SelectRelayCountsParams{
		Day:   from,
		Day_2: to,
	})
	if err != nil {
		return nil, err
	}

	var counts []types.RelayCount

	for _, dbCount := range dbCounts {
		counts = append(counts, types.RelayCount{
			AppPublicKey: dbCount.AppPublicKey,
			Day:          dbCount.Day,
			Success:      int(dbCount.Success),
			Error:        int(dbCount.Error),
			CreatedAt:    dbCount.CreatedAt,
			UpdatedAt:    dbCount.UpdatedAt,
		})
	}

	return counts, nil
}
