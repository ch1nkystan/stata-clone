package pgsql

import (
	"github.com/prosperofair/stata/pkg/types"
)

func (c *PGSQLClient) SelectBotByToken(token string) (*types.Bot, error) {
	sess := c.GetSession()

	res := &types.Bot{}

	q := `select * from bots where bot_token = ?`
	if err := sess.SelectBySql(q, token).LoadOne(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) SelectAllBots() (map[int]*types.Bot, error) {
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

func (c *PGSQLClient) SelectBotByID(id int) (*types.Bot, error) {
	sess := c.GetSession()

	res := &types.Bot{}

	q := `select * from bots where id = ?`
	if err := sess.SelectBySql(q, id).LoadOne(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) SelectBotByUsername(username string) (*types.Bot, error) {
	sess := c.GetSession()

	res := &types.Bot{}

	q := `select * from bots where bot_username = ?`
	if err := sess.SelectBySql(q, username).LoadOne(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) CreateBot(bot *types.Bot) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("bots").
		Columns(
			"api_key",
			"bot_token",
			"bot_username",
			"bot_type",
			"bid",
		).
		Record(bot).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *PGSQLClient) UpdateBotBuyerID(botToken string, bid string) error {
	sess := c.GetSession()

	if _, err := sess.Update("bots").
		Set("bid", bid).
		Where("bot_token = ?", botToken).
		Exec(); err != nil {
		return err
	}

	return nil
}
