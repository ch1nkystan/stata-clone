package pgsql

import (
	"github.com/google/uuid"
	"github.com/prosperofair/stata/pkg/types"
)

func (c *Client) SelectBotByToken(token string) (*types.Bot, error) {
	sess := c.GetSession()

	res := &types.Bot{}

	q := `select * from bots where bot_token = ? limit 1`
	if err := sess.SelectBySql(q, token).LoadOne(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectAllBots() (map[int]*types.Bot, error) {
	sess := c.GetSession()

	res := make([]*types.Bot, 0)

	q := `select * from bots`
	if _, err := sess.SelectBySql(q).Load(&res); err != nil {
		return nil, err
	}

	bots := make(map[int]*types.Bot)
	for _, bot := range res {
		bots[bot.ID] = bot
	}

	return bots, nil
}

func (c *Client) SelectBotIDsByTraceUUID(traceUUID uuid.UUID) ([]int, error) {
	sess := c.GetSession()

	res := make([]int, 0)

	q := `select id from bots where trace_uuid = ?`
	if _, err := sess.SelectBySql(q, traceUUID).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectBotByID(id int) (*types.Bot, error) {
	sess := c.GetSession()

	res := &types.Bot{}

	q := `select * from bots where id = ? limit 1`
	if err := sess.SelectBySql(q, id).LoadOne(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectBotByUsername(username string) (*types.Bot, error) {
	sess := c.GetSession()

	res := &types.Bot{}

	q := `select * from bots where bot_username = ? limit 1`
	if err := sess.SelectBySql(q, username).LoadOne(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) CreateBot(bot *types.Bot) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("bots").
		Columns(
			"api_key",
			"bot_token",
			"bot_username",
			"bot_type",
			"bid",
			"trace_uuid",
		).
		Record(bot).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBotBinding(botToken string, binding bool) error {
	sess := c.GetSession()

	q := `update bots set binding = ? where bot_token = ?`
	if _, err := sess.UpdateBySql(q, binding, botToken).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBotTraceUUID(botToken string, traceUUID uuid.UUID) error {
	sess := c.GetSession()

	q := `update bots set trace_uuid = ? where bot_token = ?`
	if _, err := sess.UpdateBySql(q, traceUUID, botToken).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateBotBuyerID(botToken string, bid string) error {
	sess := c.GetSession()

	if _, err := sess.Update("bots").
		Set("bid", bid).
		Where("bot_token = ?", botToken).
		Exec(); err != nil {
		return err
	}

	return nil
}
