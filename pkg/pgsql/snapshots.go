package pgsql

import (
	"time"

	"github.com/prosperofair/stata/pkg/types"
)

func (c *Client) CreateUsersOnlineSnapshot() error {
	sess := c.GetSession()

	q := `
		INSERT INTO snapshots (bot_id, users, snapshot, created_at)
		SELECT bot_id,
		       COUNT(telegram_id) AS online,
		       'online',
		       date_trunc('minute', now()) - INTERVAL '1 minute' * (EXTRACT(MINUTE FROM now()) % 5)
		FROM users
		WHERE messaged_at >= now() - INTERVAL '5 minutes'
		GROUP BY bot_id;
	`

	if _, err := sess.InsertBySql(q).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SelectOnlineSnapshotForInterval(botID int, start, end time.Time) (map[time.Time]int, error) {
	sess := c.GetSession()

	snapshots := make([]*types.Snapshot, 0)

	q := `select bot_id, avg(users), grid.at_interval 
		  from 	(
		  		select generate_series(min(created_at), max(created_at), interval '1 hour') as at_interval
				from snapshots
				) as grid
		  join snapshots on snapshots.created_at >= grid.at_interval
							and snapshots.created_at < grid.at_interval + interval '1 hour'
		  where bot_id = ? and grid.at_interval between ? and ?
		  group by bot_id, grid.at_interval
		  order by grid.at_interval`
	if _, err := sess.SelectBySql(q, botID, start, end).Load(&snapshots); err != nil {
		return nil, err
	}

	statistics := make(map[time.Time]int, len(snapshots))
	for _, snaphot := range snapshots {
		statistics[snaphot.CreatedAt] = snaphot.Users
	}

	return statistics, nil
}
