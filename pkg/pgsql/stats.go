package pgsql

import (
	"fmt"
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

func (c *Client) SelectUsersCountStats() ([]*UsersCountStats, error) {
	sess := c.GetSession()
	res := make([]*UsersCountStats, 0)

	q := `
		select bot_id,
       		   count(*)                                                                                           as total,
			   coalesce(sum(case when mailing_failed_attempts = 0 then 1 end), 0)                                 as success,
			   coalesce(sum(case when mailing_failed_attempts > 0 and mailing_failed_attempts < 3 then 1 end), 0) as indefinite,
			   coalesce(sum(case when mailing_failed_attempts >= 3 then 1 end), 0)                                as fail
		from users
		group by bot_id
	`
	if _, err := sess.SelectBySql(q).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectUsersCountStats: %w", err)
	}

	return res, nil
}

func (c *Client) SelectBotMailingStats(botID int) ([]*UsersCountStats, error) {
	sess := c.GetSession()
	res := make([]*UsersCountStats, 0)

	q := `
		select depot_channel_hash,
       		   count(*)                                                                                           as total,
			   coalesce(sum(case when mailing_failed_attempts = 0 then 1 end), 0)                                 as success,
			   coalesce(sum(case when mailing_failed_attempts > 0 and mailing_failed_attempts < 3 then 1 end), 0) as indefinite,
			   coalesce(sum(case when mailing_failed_attempts >= 3 then 1 end), 0)                                as fail
		from users
		where bot_id = ?
		group by depot_channel_hash
	`
	if _, err := sess.SelectBySql(q, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectBotMailingStats: %w", err)
	}

	return res, nil
}

func (c *Client) SelectBotUsersByDay(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `
		with cte as (select count(*)                                        as users_total,
							sum(case when users.seen < 1 then 1 else 0 end) as users_unique,
							date_trunc('day', users.created_at)             as by_day
             		 from users
             		 where bot_id = ?
               		   and (date(created_at) > date(?) and date(created_at) <= date(?))
					 group by by_day
					 order by by_day desc)
		select users_total,
			users_unique,
			users_unique::float4 / users_total::float4 * 100 as users_unique_rate,
			by_day
		from cte
	`
	if _, err := sess.SelectBySql(q, botID, start, end).Load(&conversions); err != nil {
		return nil, fmt.Errorf("SelectBotUsersByDay: %w", err)
	}

	for _, v := range conversions {
		res[v.ByDayDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *Client) SelectBotUsersByPeriod(botID int, period string, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `
		with cte as (select count(*)                                        as users_total,
							sum(case when users.seen < 1 then 1 else 0 end) as users_unique,
						    date_trunc(?, users.created_at)            as by_period
					 from users
					 where bot_id = ?
					   and (date(created_at) > date(?) and date(created_at) <= date(?))
					group by by_period
					order by by_period desc)
		select users_total,
			users_unique,
			users_unique::float4 / users_total::float4 * 100 as users_unique_rate,
			by_period
		from cte
	`
	if _, err := sess.SelectBySql(q, period, botID, start, end).Load(&conversions); err != nil {
		return nil, fmt.Errorf("SelectBotUsersByPeriod: %w", err)
	}

	for _, v := range conversions {
		res[v.ByPeriodDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *Client) SelectBotUsersByDeeplinks(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `
		with cte as (select count(*)                                        as users_total,
                    		sum(case when users.seen < 1 then 1 else 0 end) as users_unique,
                    		d.label                                         as label
             		 from users
                      		join deeplinks d on users.deeplink_id = d.id
             		 where users.bot_id = ?
               		   and (users.created_at > ? and users.created_at <= ?)
					 group by d.label
					 order by users_total desc)
		select users_total,
			users_unique,
			users_unique::float4 / users_total::float4 * 100 as users_unique_rate,
			label
		from cte
	`
	if _, err := sess.SelectBySql(q, botID, start, end).Load(&conversions); err != nil {
		return nil, fmt.Errorf("SelectBotUsersByDeeplinks: %w", err)
	}

	for _, v := range conversions {
		res[v.Label] = v
	}

	return res, nil
}

func (c *Client) SelectBotExpensesByDeeplinks(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `
		select sum(clicks)      as clicks,
			   sum(impressions) as impressions,
			   sum(spend)       as expense,
			   label
		from deeplinks d
			   join fbtool_accounts fa on d.label = fa.fbtool_account_name
			   join fbtool_campaigns_stats fcs on fa.fbtool_account_id = fcs.fbtool_account_id
		where d.bot_id = ?
		  and (date(fcs.date) > date(?) and date(fcs.date) <= date(?))
		group by d.label
	`
	if _, err := sess.SelectBySql(q, botID, start, end).Load(&conversions); err != nil {
		return nil, fmt.Errorf("SelectBotExpensesByDeeplinks: %w", err)
	}

	for _, v := range conversions {
		res[v.Label] = v
	}

	return res, nil
}

func (c *Client) SelectBotLeadsByDeeplinks(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)

	q := `
		with cte as (select count(distinct telegram_id)          as leads_users,
							count(*)                             as leads_total,
							coalesce(sum(t.price * t.amount), 0) as income,
							d.label
             		 from users u
                      		join deeplinks d on d.id = u.deeplink_id
                      		full outer join transactions t on u.id = t.user_id
					 where u.bot_id = ?
					   and (date(u.deposited_at) > date(?) and date(u.deposited_at) <= date(?))
					   and (date(t.created_at) > date(?) and date(t.created_at) <= date(?) or t.created_at is null)
					   and deposited = true
					 group by d.label)
		select leads_users,
			   leads_total,
			   leads_total::float4 / leads_users::float4 as leads_per_user,
			   income,
			   label
		from cte
	`
	if _, err := sess.SelectBySql(q, botID, start, end, start, end).Load(&conversions); err != nil {
		return nil, fmt.Errorf("SelectBotLeadsByDeeplinks: %w", err)
	}

	for _, v := range conversions {
		res[v.Label] = v
	}

	return res, nil
}

func (c *Client) SelectBotLeadsByPeriod(botID int, period string, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `
		with cte as (select count(distinct telegram_id)            as leads_users,
							count(*)                               as leads_total,
							coalesce(sum(t.price * t.amount), 0)   as income,
							date_trunc(?, users.deposited_at) 	   as by_period
					 from users
							full outer join transactions t on users.id = t.user_id
					 where bot_id = ?
					   and (date(users.deposited_at) > date(?) and date(users.deposited_at) <= date(?))
					   and (date(t.created_at) > date(?) and date(t.created_at) <= date(?) or t.created_at is null)
					   and deposited = true
					 group by by_period
					 order by by_period desc)
		select leads_users,
			   leads_total,
			   leads_total::float4 / leads_users::float4 as leads_per_user,
			   income,
			   by_period
		from cte
	`
	if _, err := sess.SelectBySql(q, period, botID, start, end, start, end).Load(&conversions); err != nil {
		return nil, fmt.Errorf("SelectBotLeadsByPeriod: %w", err)
	}

	for _, v := range conversions {
		res[v.ByPeriodDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *Client) SelectBotLeadsByDay(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	conversions := make([]*types.ConversionRow, 0)
	q := `
		with cte as (select count(distinct telegram_id)           as leads_users,
							count(*)                              as leads_total,
							coalesce(sum(t.price * t.amount), 0)  as income,
							date_trunc('day', users.deposited_at) as by_day
					 from users
							full outer join transactions t on users.id = t.user_id
					 where bot_id = ?
					   and (date(users.deposited_at) > date(?) and date(users.deposited_at) <= date(?))
					   and (date(t.created_at) > date(?) and date(t.created_at) <= date(?) or t.created_at is null)
					   and users.deposited = true
					 group by by_day
					 order by by_day desc)
		select leads_users,
			   leads_total,
			   leads_total::float4 / leads_users::float4 as leads_per_user,
			   income,
			   by_day
		from cte
	`
	if _, err := sess.SelectBySql(q, botID, start, end, start, end).Load(&conversions); err != nil {
		return nil, fmt.Errorf("SelectBotLeadsByDay: %w", err)
	}

	for _, v := range conversions {
		res[v.ByDayDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *Client) SelectBotExpensesByPeriod(botID int, period string, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	expenses := make([]*types.ConversionRow, 0)
	q := `
		select sum(clicks)                  as clicks,
			   sum(impressions)             as impressions,
			   sum(spend)                   as expense,
			   date_trunc(?, fcs.date) 		as by_period
		from deeplinks d
			   join fbtool_accounts fa on d.label = fa.fbtool_account_name
			   join fbtool_campaigns_stats fcs on fa.fbtool_account_id = fcs.fbtool_account_id
		where d.bot_id = ?
		  and (date(fcs.date) > date(?) and date(fcs.date) <= date(?))
		group by by_period
	`
	if _, err := sess.SelectBySql(q, period, botID, start, end).Load(&expenses); err != nil {
		return nil, fmt.Errorf("SelectBotExpensesByPeriod: %w", err)
	}

	for _, v := range expenses {
		res[v.ByPeriodDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *Client) SelectBotExpensesByDay(botID int, start, end time.Time) (map[string]*types.ConversionRow, error) {
	sess := c.GetSession()
	res := make(map[string]*types.ConversionRow, 0)
	expenses := make([]*types.ConversionRow, 0)
	q := `
		select sum(clicks)                 as clicks,
			   sum(impressions)            as impressions,
			   sum(spend)                  as expense,
			   date_trunc('day', fcs.date) as by_day
		from deeplinks d
         	   join fbtool_accounts fa on d.label = fa.fbtool_account_name
         	   join fbtool_campaigns_stats fcs on fa.fbtool_account_id = fcs.fbtool_account_id
		where d.bot_id = ?
		  and (date(fcs.date) > date(?) and date(fcs.date) <= date(?))
		group by by_day
	`
	if _, err := sess.SelectBySql(q, botID, start, end).Load(&expenses); err != nil {
		return nil, fmt.Errorf("SelectBotExpensesByDay: %w", err)
	}

	for _, v := range expenses {
		res[v.ByDayDB.Format(time.DateOnly)] = v
	}

	return res, nil
}

func (c *Client) SelectDepositsByBotID(botID int, start, end time.Time) ([]*types.DepositRow, error) {
	sess := c.GetSession()
	res := make([]*types.DepositRow, 0)

	q := `
		select tx.id as id,
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
		  and (date(tx.created_at) > date(?) and date(tx.created_at) <= date(?))
		order by tx.created_at desc
	`
	if _, err := sess.SelectBySql(q, botID, start, end).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectDepositsByBotID: %w", err)
	}

	return res, nil
}

func (c *Client) SelectUsersMetric(botID int, start, end, startPrev, endPrev time.Time) (*types.MetricRow, error) {
	sess := c.GetSession()
	res := &types.MetricRow{}

	q := `
		with cte as (select count(*)                                                             						as total,
                    		sum(case when date(created_at) > date(?) and date(created_at) <= date(?) then 1 else 0 end) as current_period,
                    		sum(case when date(created_at) > date(?) and date(created_at) <= date(?) then 1 else 0 end) as last_period
					 from users
					 where mailing_state = 'ready' and bot_id = ?)
		select total                                                                                                as all_time,
			   current_period                                                                                       as period,
			   last_period 																							as last_period,
			   case when last_period = 0 then 100 else float4(current_period) / float4(last_period) * 100 - 100 end as diff
		from cte
	`
	if _, err := sess.SelectBySql(q, start, end, startPrev, endPrev, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectUsersMetric: %w", err)
	}

	return res, nil
}

func (c *Client) SelectUsersReferralsMetric(botID int, start, end, startPrev, endPrev time.Time) (*types.MetricRow, error) {
	sess := c.GetSession()
	res := &types.MetricRow{}

	q := `
		with cte2 as (with cte as (select count(*)                                                            					   				   						   as total,
										  sum(case when date(u.created_at) > date(?) and date(u.created_at) <= date(?) then 1 else 0 end) 					   			   as cp_total,
										  sum(case when date(u.created_at) > date(?) and date(u.created_at) <= date(?) then 1 else 0 end) 					   			   as lp_total,
										  sum(case when dl.referral_telegram_id != 0 then 1 else 0 end) 									       						   as refs,
										  sum(case when dl.referral_telegram_id != 0 and date(u.created_at) > date(?) and date(u.created_at) <= date(?) then 1 else 0 end) as cp_refs,
										  sum(case when dl.referral_telegram_id != 0 and date(u.created_at) > date(?) and date(u.created_at) <= date(?) then 1 else 0 end) as lp_refs
								   from users u
								   		  left join deeplinks dl on u.deeplink_id = dl.id
								   where u.bot_id = ?)

					  select case when coalesce(total, 0) = 0 then 0 else refs::float4 / total::float4 * 100 end          as all_time,
                      		 case when coalesce(cp_total, 0) = 0 then 0 else cp_refs::float4 / cp_total::float4 * 100 end as period,
                      		 case when coalesce(lp_total, 0) = 0 then 0 else lp_refs::float4 / lp_total::float4 * 100 end as last_period
                      from cte)

		select all_time, period, last_period, period - last_period as diff
		from cte2
	`
	if _, err := sess.SelectBySql(q, start, end, startPrev, endPrev, start, end, startPrev, endPrev, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectUsersReferralsMetric: %w", err)
	}

	return res, nil
}

func (c *Client) SelectLeadsMetric(botID int, start, end, startPrev, endPrev time.Time) (*types.MetricRow, error) {
	sess := c.GetSession()
	res := &types.MetricRow{}

	q := `
		with cte as (select count(*)                                                             						as total,
							sum(case when date(created_at) > date(?) and date(created_at) <= date(?) then 1 else 0 end) as current_period,
							sum(case when date(created_at) > date(?) and date(created_at) <= date(?) then 1 else 0 end) as last_period
					 from users
					 where deposited = true and bot_id = ?)
		select total                                                                                                as all_time,
			   current_period                                                                                       as period,
			   last_period 																							as last_period,
			   case when last_period = 0 then 100 else float4(current_period) / float4(last_period) * 100 - 100 end as diff
		from cte
	`
	if _, err := sess.SelectBySql(q, start, end, startPrev, endPrev, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectLeadsMetric: %w", err)
	}

	return res, nil
}

func (c *Client) SelectExpenseMetric(botID int, start, end, startPrev, endPrev time.Time) (*types.MetricRow, error) {
	sess := c.GetSession()
	res := &types.MetricRow{}

	q := `
		with cte as (select sum(spend)                                                          						as total,
                    		sum(case when date(fcs.date) > date(?) and date(fcs.date) <= date(?) then spend else 0 end) as current_period,
                    		sum(case when date(fcs.date) > date(?) and date(fcs.date) <= date(?) then spend else 0 end) as last_period
             		 from deeplinks d
                      		join fbtool_accounts fa on d.label = fa.fbtool_account_name
                      		join fbtool_campaigns_stats fcs on fa.fbtool_account_id = fcs.fbtool_account_id
             		 where d.bot_id = ?)

		select coalesce(total, 0)                                                    as all_time,
			   coalesce(current_period, 0)                                           as period,
			   coalesce(last_period, 0)                                              as last_period,
			   case 
			   	   when coalesce(last_period, 0) = 0 then 0 
				   else float4(current_period) / float4(last_period) * 100 - 100 end as diff
		from cte
	`
	if _, err := sess.SelectBySql(q, start, end, startPrev, endPrev, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectExpenseMetric: %w", err)
	}

	return res, nil
}

func (c *Client) SelectIncomeMetric(botID int, start, end, startPrev, endPrev time.Time) (*types.MetricRow, error) {
	sess := c.GetSession()
	res := &types.MetricRow{}

	q := `
		with cte as (select sum(t.amount * t.price)                                                             							 as total,
							sum(case when date(t.created_at) > date(?) and date(t.created_at) <= date(?) then t.amount * t.price else 0 end) as current_period,
							sum(case when date(t.created_at) > date(?) and date(t.created_at) <= date(?) then t.amount * t.price else 0 end) as last_period
					from users join public.transactions t on users.id = t.user_id
					where deposited = true
					  and bot_id = ?)
		select coalesce(total, 0)                                                    as all_time,
			   coalesce(current_period, 0)                                           as period,
			   coalesce(last_period, 0)                                              as last_period,
			   case
				   when coalesce(last_period, 0) = 0 then 0
				   else float4(current_period) / float4(last_period) * 100 - 100 end as diff
		from cte
	`
	if _, err := sess.SelectBySql(q, start, end, startPrev, endPrev, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectIncomeMetric: %w", err)
	}

	return res, nil
}

func (c *Client) SelectClicksMetric(botID int, start, end, startPrev, endPrev time.Time) (*types.MetricRow, error) {
	sess := c.GetSession()
	res := &types.MetricRow{}

	q := `
		with cte as (select sum(clicks)                                                          as total,
                    		sum(case when date(fcs.date) > date(?) and date(fcs.date) <= date(?) then clicks else 0 end) as current_period,
                    		sum(case when date(fcs.date) > date(?) and date(fcs.date) <= date(?) then clicks else 0 end) as last_period
             		 from deeplinks d
                      		  join fbtool_accounts fa on d.label = fa.fbtool_account_name
                      		  join fbtool_campaigns_stats fcs on fa.fbtool_account_id = fcs.fbtool_account_id
             		 where d.bot_id = ?)

		select coalesce(total, 0)                                                    as all_time,
			   coalesce(current_period, 0)                                           as period,
			   coalesce(last_period, 0)                                              as last_period,
			   case 
			   	   when coalesce(last_period, 0) = 0 then 0 
				   else float4(current_period) / float4(last_period) * 100 - 100 end as diff
		from cte
	`
	if _, err := sess.SelectBySql(q, start, end, startPrev, endPrev, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectClicksMetric: %w", err)
	}

	return res, nil
}

func (c *Client) SelectUsersUniqueMetric(botID int, start, end, startPrev, endPrev time.Time) (*types.MetricRow, error) {
	sess := c.GetSession()
	res := &types.MetricRow{}

	q := `
		with cte2 as (with cte as (select count(*)                                                                         as total,
										  sum(case when seen < 1 then 1 else 0 end)                                        as uq,
										  sum(case when date(created_at) > date(?) and date(created_at) <= date(?) then 1 else 0 end)              as cp_total,
										  sum(case when date(created_at) > date(?) and date(created_at) <= date(?) and seen < 1 then 1 else 0 end) as cp_uq,
										  sum(case when date(created_at) > date(?) and date(created_at) <= date(?) then 1 else 0 end)              as lp_total,
										  sum(case when date(created_at) > date(?) and date(created_at) <= date(?) and seen < 1 then 1 else 0 end) as lp_uq
								   from users
								   where bot_id = ?)
					  select case when total = 0 then 0 else uq::float4 / total::float4 * 100 end          as all_time,
							 case when cp_total = 0 then 0 else cp_uq::float4 / cp_total::float4 * 100 end as period,
							 case when lp_total = 0 then 0 else lp_uq::float4 / lp_total::float4 * 100 end as last_period
					  from cte)

		select all_time, period, last_period, period - last_period as diff 
		from cte2
	`
	if _, err := sess.SelectBySql(q, start, end, start, end, startPrev, endPrev, startPrev, endPrev, botID).Load(&res); err != nil {
		return nil, fmt.Errorf("SelectUsersUniqueMetric: %w", err)
	}

	return res, nil
}
