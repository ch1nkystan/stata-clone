package pgsql

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/prosperofair/stata/pkg/types"
)

func (c *Client) CreateFBToolAccount(record *types.FBToolAccount) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("fbtool_accounts").
		Columns(
			"token_id",
			"fbtool_account_id",
			"fbtool_account_name",
			"fetched_at",
		).Record(record).Exec(); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return ErrAlreadyExists
			}
		}
		return fmt.Errorf("failed to create fbtool account: %w", err)
	}

	return nil
}

func (c *Client) SelectUnfetchedFBToolTokens() ([]*types.FBToolToken, error) {
	sess := c.GetSession()

	q := `select * from fbtool_tokens where active = true and fetched_at < now() - interval '6 hour' order by random() limit 1;`

	res := make([]*types.FBToolToken, 0)
	if _, err := sess.SelectBySql(q).Load(&res); err != nil {
		return nil, fmt.Errorf("failed to select unfetched fbtool tokens: %w", err)
	}

	return res, nil
}

func (c *Client) SelectUnfetchedFBToolAccounts(tokenID int) ([]*types.FBToolAccount, error) {
	sess := c.GetSession()

	q := `select * from fbtool_accounts where active = true and token_id = ? and fetched_at < now() - interval '1 day' and fetch_duration < 25`

	res := make([]*types.FBToolAccount, 0)
	if _, err := sess.SelectBySql(q, tokenID).Load(&res); err != nil {
		return nil, fmt.Errorf("failed to select unfetched fbtool accounts: %w", err)
	}

	return res, nil
}

// type FBToolCampaignStatRecord struct {
// 	CampaignID string    `db:"campaign_id"`
// 	Date       time.Time `db:"date"`
// }

// func (c *Client) SelectFBToolCampaignStatsRecordsByDay(accountID, daysToFetch int) ([]*FBToolCampaignStatRecord, error) {
// 	sess := c.GetSession()

// 	q := `
// 		select count(*), campaign_id, date
// 		from fbtool_campaigns_stats
// 		where date >= now() - interval '? days'
//   		  and fbtool_account_id = ?
// 		group by campaign_id, date
// 	`
// 	res := make([]*FBToolCampaignStatRecord, 0)
// 	if _, err := sess.SelectBySql(q, daysToFetch, accountID).Load(&res); err != nil {
// 		return nil, fmt.Errorf("failed to select fbtool campaign stats records by day: %w", err)
// 	}

// 	return res, nil
// }

func (c *Client) UpdateFBToolTokenFetchedAt(id int) error {
	sess := c.GetSession()

	q := `update fbtool_tokens set fetched_at = now() where id = ?`
	if _, err := sess.UpdateBySql(q, id).Exec(); err != nil {
		return fmt.Errorf("failed to update fbtool token fetched_at: %w", err)
	}

	return nil
}

func (c *Client) UpdateFBToolTokenDaysToFetch(id int) error {
	sess := c.GetSession()

	q := `update fbtool_tokens set days_to_fetch = 2 where id = ?`
	if _, err := sess.UpdateBySql(q, id).Exec(); err != nil {
		return fmt.Errorf("failed to update fbtool token days_to_fetch: %w", err)
	}

	return nil
}

func (c *Client) UpdateFBToolAccountFetchedAt(aid int, fetched bool, fd int) error {
	sess := c.GetSession()

	q := `update fbtool_accounts set fetched_at = now(), fetched = ?, fetch_duration = ? where fbtool_account_id = ?`
	if _, err := sess.UpdateBySql(q, fetched, fd, aid).Exec(); err != nil {
		return fmt.Errorf("failed to update fbtool account fetched_at: %w", err)
	}

	return nil
}

// func (c *Client) DisableFBToolAccount(aid int) error {
// 	sess := c.GetSession()

// 	q := `update fbtool_accounts set active = false where fbtool_account_id = ?`
// 	if _, err := sess.UpdateBySql(q, aid).Exec(); err != nil {
// 		return fmt.Errorf("failed to disable fbtool account: %w", err)
// 	}

// 	return nil
// }

func (c *Client) CreateFBToolCampaignStat(record *types.FBToolCampaignStat) error {
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
		return fmt.Errorf("failed to create fbtool campaign stat: %w", err)
	}

	return nil
}
