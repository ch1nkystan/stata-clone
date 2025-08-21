package pgsql

import (
	"time"

	"github.com/prosperofair/stata/pkg/types"
)

const DateFormatDay = "2006-01-02"

func (c *Client) SelectLeadsByCampaign(token string, start, end time.Time) ([]*types.LeadsByCampaignRow, error) {
	sess := c.GetSession()

	res := make([]*types.LeadsByCampaignRow, 0)

	q := `SELECT d.label,
       count(*)                                         AS users_total,
       sum(case when users.seen < 1 then 1 else 0 end)  as users_unique,
       sum(case when users.deposited then 1 else 0 end) as users_deposited,
       sum(users.deposits_sum)                          as deposits_sum,
       sum(users.deposits_total)                        as deposits_total,
       date_trunc('day', users.created_at)              AS day
FROM users
         JOIN bots b ON b.id = users.bot_id
         JOIN deeplinks d ON d.id = users.deeplink_id
WHERE b.bot_token = ?
  AND date_trunc('day', users.created_at) BETWEEN date ?
    AND date ?
GROUP BY d.label, day
order by users_total desc;
`
	if _, err := sess.SelectBySql(q, token,
		start.Format(DateFormatDay),
		end.Format(DateFormatDay)).Load(&res); err != nil {
		return nil, err
	}

	for i := range res {
		if float64(res[i].UsersTotal) > 0 {
			res[i].UsersUniqueRate = float64(res[i].UsersUnique) / float64(res[i].UsersTotal)
		}

		if float64(res[i].UsersTotal) > 0 {
			res[i].DepositsPerUser = float64(res[i].DepositsTotal) / float64(res[i].UsersTotal)
		}
	}

	return res, nil
}
