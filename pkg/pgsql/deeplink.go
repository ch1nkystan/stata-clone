package pgsql

import (
	"github.com/prosperofair/stata/pkg/types"
)

func (c *Client) CreateDeeplink(deeplink *types.Deeplink) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("deeplinks").
		Columns(
			"bot_id",
			"referral_telegram_id",
			"hash",
			"label",
		).Record(deeplink).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SelectBotDeeplinks(botID int) ([]*types.Deeplink, error) {
	sess := c.GetSession()

	res := make([]*types.Deeplink, 0)

	q := `select * from deeplinks where bot_id = ? order by id desc`
	if _, err := sess.SelectBySql(q, botID).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) UpdateDeeplinkLabel(botID int, hash, label string) error {
	sess := c.GetSession()

	q := `update deeplinks set label = ? where bot_id = ? and hash = ?`
	if _, err := sess.UpdateBySql(q, label, botID, hash).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SelectBotDeeplinksByHash(botID int, hash string) ([]*types.Deeplink, error) {
	sess := c.GetSession()

	res := make([]*types.Deeplink, 0)

	q := `select * from deeplinks where bot_id = ? and hash = ? order by id desc limit 1`
	if _, err := sess.SelectBySql(q, botID, hash).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectBotDeeplinksByReferralID(botID int, referralID int64) ([]*types.Deeplink, error) {
	sess := c.GetSession()

	res := make([]*types.Deeplink, 0)

	q := `select * from deeplinks where bot_id = ? and referral_telegram_id = ? order by id desc`
	if _, err := sess.SelectBySql(q, botID, referralID).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SelectBotDeeplinksByLabel(botID int, label string) ([]*types.Deeplink, error) {
	sess := c.GetSession()

	res := make([]*types.Deeplink, 0)

	q := `select * from deeplinks where bot_id = ? and label = ? order by id desc`
	if _, err := sess.SelectBySql(q, botID, label).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}
