package pgsql

import (
	"time"

	"github.com/prosperofair/stata/pkg/types"
)

type UsersCountStats struct {
	BotID            int    `db:"bot_id"`
	DepotChannelHash string `db:"depot_channel_hash"`
	Total            int    `db:"total"`
	Success          int    `db:"success"`
	Indefinite       int    `db:"indefinite"`
	Fail             int    `db:"fail"`
}

func (c *PGSQLClient) SelectUsersCountStats() ([]*UsersCountStats, error) {
	sess := c.GetSession()
	res := make([]*UsersCountStats, 0)

	q := `select bot_id,
       count(*)                                                                                           as total,
       coalesce(sum(case when mailing_failed_attempts = 0 then 1 end), 0)                                 as success,
       coalesce(sum(case when mailing_failed_attempts > 0 and mailing_failed_attempts < 3 then 1 end), 0) as indefinite,
       coalesce(sum(case when mailing_failed_attempts >= 3 then 1 end), 0)                                as fail
from users
group by bot_id;`

	if _, err := sess.SelectBySql(q).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) SelectBotMailingStats(botID int) ([]*UsersCountStats, error) {
	sess := c.GetSession()
	res := make([]*UsersCountStats, 0)

	q := `select depot_channel_hash,
       count(*)                                                                                           as total,
       coalesce(sum(case when mailing_failed_attempts = 0 then 1 end), 0)                                 as success,
       coalesce(sum(case when mailing_failed_attempts > 0 and mailing_failed_attempts < 3 then 1 end), 0) as indefinite,
       coalesce(sum(case when mailing_failed_attempts >= 3 then 1 end), 0)                                as fail
from users
where bot_id = ?
group by depot_channel_hash;`

	if _, err := sess.SelectBySql(q, botID).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) SelectBotUsersByDay(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `with cte as (select count(*)                                        as users_total,
                    sum(case when users.seen < 1 then 1 else 0 end) as users_unique,
                    date_trunc('day', users.created_at)             as by_day
             from users
             where bot_id = ?
               and created_at >= ?
               and created_at <= ?
             group by by_day
             order by by_day desc)
select users_total,
       users_unique,
       users_unique::float4 / users_total::float4 * 100 as users_unique_rate,
       by_day
from cte;`

	if _, err := sess.SelectBySql(q, botID, start, end).Load(&conversions); err != nil {
		return nil, err
	}

	withDeeplinks, err := c.SelectBotUsersWithDeeplinksByDay(botID, start, end)
	if err != nil {
		return nil, err
	}

	for _, v := range conversions {
		for _, d := range withDeeplinks {
			if v.Deeplinks == nil {
				v.Deeplinks = make([]*types.ConversionRow, 0)
			}

			if v.ByDayDB.Format(time.DateOnly) == d.ByDayDB.Format(time.DateOnly) {
				v.Deeplinks = append(v.Deeplinks, d)
			}
		}

		res[v.ByDayDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *PGSQLClient) SelectBotUsersWithDeeplinksByDay(botID int, start, end time.Time) ([]*types.ConversionRow, error) {
	sess := c.GetSession()

	conversions := make([]*types.ConversionRow, 0)
	q := `with cte as (select count(*)                                        as users_total,
                    sum(case when users.seen < 1 then 1 else 0 end) as users_unique,
                    d.label                                         as label,
                    date_trunc('day', users.created_at)             as by_day
             from users
                      join deeplinks d on users.deeplink_id = d.id
             where users.bot_id = ?
               and users.created_at >= ?
               and users.created_at <= ?
             group by d.label, by_day
             order by by_day desc)
select users_total,
       users_unique,
       users_unique::float4 / users_total::float4 * 100 as users_unique_rate,
       label,
       by_day
from cte;`

	if _, err := sess.SelectBySql(q, botID, start, end).Load(&conversions); err != nil {
		return nil, err
	}

	return conversions, nil
}

func (c *PGSQLClient) SelectBotLeadsByDay(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `with cte as (select count(distinct telegram_id)           as leads_users,
                    count(*)                              as leads_total,
                    coalesce(sum(t.price * t.amount), 0)  as profit,
                    date_trunc('day', users.deposited_at) as by_day
             from users
                      full outer join transactions t on users.id = t.user_id
             where bot_id = ?
               and users.deposited_at >= ?
               and users.deposited_at <= ?
               and (t.created_at >= ? and t.created_at <= ? or t.created_at is null)
               and deposited = true
             group by by_day
             order by by_day desc)
select leads_users,
       leads_total,
       leads_total::float4 / leads_users::float4 as leads_per_user,
	   profit,
       by_day
from cte;`

	if _, err := sess.SelectBySql(q, botID, start, end, start, end).Load(&conversions); err != nil {
		return nil, err
	}

	for _, v := range conversions {
		res[v.ByDayDB.Format(time.DateOnly)] = v
	}

	withDeeplinks, err := c.SelectBotLeadsWithDeeplinksByDay(botID, start, end)
	if err != nil {
		return nil, err
	}

	for _, v := range conversions {
		for _, d := range withDeeplinks {
			if v.DeeplinksLeads == nil {
				v.DeeplinksLeads = make(map[string]*types.ConversionRow, 0)
			}

			if v.ByDayDB.Format(time.DateOnly) == d.ByDayDB.Format(time.DateOnly) {
				v.DeeplinksLeads[d.Label] = d
			}
		}

		res[v.ByDayDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *PGSQLClient) SelectBotLeadsWithDeeplinksByDay(botID int, start, end time.Time) ([]*types.ConversionRow, error) {
	sess := c.GetSession()

	conversions := make([]*types.ConversionRow, 0)
	q := `with cte as (select count(distinct telegram_id)           as leads_users,
                    count(*)                              as leads_total,
                    d.label                               as label,
                    coalesce(sum(t.price * t.amount), 0)  as profit,
                    date_trunc('day', users.deposited_at) as by_day
             from users
                      join deeplinks d on users.deeplink_id = d.id
                      full outer join transactions t on users.id = t.user_id
             where users.bot_id = ?
               and users.deposited_at >= ?
               and users.deposited_at <= ?
               and (t.created_at >= ? and t.created_at <= ? or t.created_at is null)
               and deposited = true
             group by d.label, by_day
             order by by_day desc)
select leads_users,
       leads_total,
       leads_total::float4 / leads_users::float4 as leads_per_user,
       profit,
       label,
       by_day
from cte;`

	if _, err := sess.SelectBySql(q, botID, start, end, start, end).Load(&conversions); err != nil {
		return nil, err
	}

	return conversions, nil
}

func (c *PGSQLClient) SelectDepositsByBotID(botID int, start, end time.Time) ([]*types.DepositRow, error) {
	sess := c.GetSession()
	res := make([]*types.DepositRow, 0)

	q := `select tx.id as id,
	   tx.user_id 	        as user_id,
       tx.tx_key            as hash,
       dl.label             as deeplink,
       tx.blockchain        as blockchain,
       tx.amount * tx.price as amount,
	   tx.created_at        as date
from transactions tx
         join public.users u on u.id = tx.user_id
		 join deeplinks dl on u.deeplink_id = dl.id
where u.bot_id = ?
  and u.deposited = true
  and tx.created_at >= ?
  and tx.created_at <= ?
order by tx.created_at desc;`

	if _, err := sess.SelectBySql(q, botID, start, end).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}
