package pgsql

import (
	"time"

	"github.com/prosperofair/stata/pkg/types"
)

func (c *Client) SelectLeadsByCampaign(token string, start, end time.Time) ([]*types.LeadByCampaignRow, error) {
	sess := c.GetSession()

	res := make([]*types.LeadByCampaignRow, 0, 0)

	q := `with cte4 as (with cte as (select d.label,
                                  			count(*)                                         as total,
											sum(case when users.seen < 1 then 1 else 0 end)  as uq,
											sum(case when users.deposited then 1 else 0 end) as deposited,
											date_trunc('day', users.created_at)              as day,
											sum(users.deposits_sum)                          as deposits_sum,
											sum(users.deposits_total)                        as deposits_total
                           			from users
                                    join bots b on b.id = users.bot_id
                                    join deeplinks d on d.id = users.deeplink_id
                               			and b.bot_token = ?
									group by d.label, day
									order by day desc),
                   			cte2 as (select fba.fbtool_account_name     as label,
                                   			date_trunc('day', fcs.date) as day,
                                   			sum(spend)                  as total_spend
                            		from fbtool_accounts fba
                                	join public.fbtool_campaigns_stats fcs
                                    	on fba.fbtool_account_id = fcs.fbtool_account_id
									group by fba.fbtool_account_name, day
									order by day desc),
                   			cte3 as (select date_id
                            		from generate_series(?, ?,'1 day'::interval) as date_id)
        	select cte.label                                      	  as label,
                	sum(total)                                        as total,
                	sum(uq)                                           as unique,
                    sum(uq)::float4 / sum(total)::float4 * 100        as uniqueness_ratio,
                    sum(deposited)                                    as deposited,
                    float4(sum(deposited)) / float4(sum(total)) * 100 as deposit_rate,
                    coalesce(sum(total_spend), 0)                     as total_spend,
                    sum(deposits_sum)                                 as deposits_sum,
                    sum(deposits_total)                               as deposits_total,
                    coalesce(sum(total_spend), 0) / sum(total) 		  as cpu,
                    case
                        when sum(deposited) > 0 then sum(deposits_total) / sum(deposited)
                        else 0 end                                    as deposits_per_person,
                    case
                        when sum(deposited) > 0 then coalesce(sum(total_spend), 0) / sum(deposited)
                        else 0 end                                    as deposit_price,
                    coalesce(sum(deposits_sum) - sum(total_spend), 0) as pnl
            from cte3
            join cte on cte3.date_id = cte.day
            full outer join cte2 on cte3.date_id = cte2.day 
				and lower(replace(cte2.label, ' ', '')) = lower(replace(cte.label, ' ', ''))
            group by cte.label)
	select *
	from cte4
	where label is not null
	order by total desc;
`
	if _, err := sess.SelectBySql(q, token, start, end).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}
