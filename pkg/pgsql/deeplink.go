package pgsql

import (
	"fmt"
	"github.com/prosperofair/pkg/log"
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

func (c *Client) CreateDeeplinksIgnoreInsertErrorsReturningAll(deeplink []types.Deeplink) ([]types.Deeplink, error) {
	sess := c.GetSession()

	stmt := sess.InsertInto("deeplinks").Columns(
		"bot_id",
		"referral_telegram_id",
		"hash",
		"label",
		"active",
	)

	for _, deeplink := range deeplink {
		stmt = stmt.Record(deeplink)
	}

	var deeplinks []types.Deeplink
	err := stmt.Returning(
		"id",
		"bot_id",
		"referral_telegram_id",
		"active",
		"hash",
		"label",
		"created_at",
		"updated_at",
	).Load(&deeplinks)

	return deeplinks, err
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

func (c *Client) SelectDeeplinksByHashes(hashes []string) ([]*types.Deeplink, error) {
	sess := c.GetSession()
	var res []*types.Deeplink

	_, err := sess.Select("*").
		From("deeplinks").
		Where("hash IN ?", hashes).
		OrderDesc("id").
		Load(&res)

	return res, err
}

func (c *Client) SelectBotDeeplinksByReferralID(botID int, referralID int64, limit uint64) ([]*types.Deeplink, error) {
	sess := c.GetSession()

	res := make([]*types.Deeplink, 0)

	q := `select * from deeplinks where bot_id = ? and referral_telegram_id = ? order by id desc`

	stmt := sess.SelectBySql(q, botID, referralID)

	if limit != 0 {
		log.Info(fmt.Sprintf("limit: %d", limit))
		stmt = stmt.Limit(limit)
	}

	if _, err := stmt.Load(&res); err != nil {
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

func (c *Client) SelectPixelLink(inviteUUID string) (*types.PixelLink, error) {
	sess := c.GetSession()

	res := make([]*types.PixelLink, 0)

	q := `select * from pixel_links where invite_uuid = ? order by id desc limit 1`
	if _, err := sess.SelectBySql(q, inviteUUID).Load(&res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}
