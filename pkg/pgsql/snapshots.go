package pgsql

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
