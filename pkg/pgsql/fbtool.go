package pgsql

import (
	"fmt"
	"time"

	"github.com/prosperofair/stata/pkg/types"
)

func (c *PGSQLClient) CreateFBToolAccount(record *types.FBToolAccount) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("fbtool_accounts").
		Columns(
			"token_id",
			"fbtool_account_id",
			"fbtool_account_name",
		).Record(record).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *PGSQLClient) SelectUnfetchedFBToolTokens() ([]*types.FBToolToken, error) {
	sess := c.GetSession()

	res := make([]*types.FBToolToken, 0)
	q := `select * from fbtool_tokens where active = true and fetched_at < now() - interval '6 hour'`
	if _, err := sess.SelectBySql(q).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) SelectUnfetchedFBToolAccounts(tokenID int) ([]*types.FBToolAccount, error) {
	sess := c.GetSession()

	res := make([]*types.FBToolAccount, 0)
	q := `select * from fbtool_accounts where active = true and token_id = ? and fetched_at < now() - interval '1 day'`
	if _, err := sess.SelectBySql(q, tokenID).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

type FBToolCampaignStatRecord struct {
	CampaignID string    `db:"campaign_id"`
	Date       time.Time `db:"date"`
}

func (c *PGSQLClient) SelectFBToolCampaignStatsRecordsByDay(accountID, daysToFetch int) ([]*FBToolCampaignStatRecord, error) {
	sess := c.GetSession()

	res := make([]*FBToolCampaignStatRecord, 0)
	qtpl := `select count(*), campaign_id, date
	from fbtool_campaigns_stats
	where date >= now() - interval '%d days'
	  and fbtool_account_id = ?
	group by campaign_id, date`

	q := fmt.Sprintf(qtpl, daysToFetch)

	if _, err := sess.SelectBySql(q, accountID).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) SelectFBToolAccountsByAccountID(aid int) ([]*types.FBToolAccount, error) {
	sess := c.GetSession()

	res := make([]*types.FBToolAccount, 0)
	q := `select * from fbtool_accounts where fbtool_account_id = ? limit 1`
	if _, err := sess.SelectBySql(q, aid).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) UpdateFBToolTokenFetchedAt(id int) error {
	sess := c.GetSession()

	q := `update fbtool_tokens set fetched_at = now() where id = ?`
	if _, err := sess.UpdateBySql(q, id).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *PGSQLClient) UpdateFBToolAccountFetchedAt(aid int) error {
	sess := c.GetSession()

	q := `update fbtool_accounts set fetched_at = now() where fbtool_account_id = ?`
	if _, err := sess.UpdateBySql(q, aid).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *PGSQLClient) DisableFBToolAccount(aid int) error {
	sess := c.GetSession()

	q := `update fbtool_accounts set active = false where fbtool_account_id = ?`
	if _, err := sess.UpdateBySql(q, aid).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *PGSQLClient) CreateFBToolCampaignStat(record *types.FbToolCampaignStat) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("fbtool_campaigns_stats").
		Columns(
			"fbtool_account_id",
			"campaign_name",
			"campaign_id",
			"status",
			"effective_status",
			"impressions",
			"clicks",
			"spend",
			"date",
		).Record(record).Exec(); err != nil {
		return err
	}

	return nil
}
