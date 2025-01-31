package pgsql

import (
	"fmt"
	"time"
)

func (c *Client) CreateOnlineHistory(start, end time.Time) error {
	sess := c.GetSession()

	q := `
		insert into online_history (bot_id, interval_start, interval_end, active_users_count)
		select 
			bot_id,
			? as interval_start,
			? as interval_end,
			count(distinct telegram_id) as active_users_count
		from 
			users
		where 
			messaged_at BETWEEN ? AND ?
		group by 
			bot_id;
	`

	if _, err := sess.InsertBySql(q, start, end, start, end).Exec(); err != nil {
		return fmt.Errorf("SelectUsersUniqueMetric: %w", err)
	}

	return nil
}
