package pgsql

import "github.com/prosperofair/stata/pkg/types"

func (c *Client) CreateEventLog(el *types.EventsLog) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("events_log").
		Columns(
			"event_type",
			"reporter_telegram_id",
			"user_id",
		).Record(el).Exec(); err != nil {
		return err
	}

	return nil
}
