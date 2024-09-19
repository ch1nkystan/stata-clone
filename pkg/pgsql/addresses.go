package pgsql

import (
	"fmt"

	"github.com/prosperofair/stata/pkg/types"
)

func (c *Client) CreateAddress(address *types.Address) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("addresses").
		Columns(
			"blockchain",
			"address_key",
			"address",
			"bid",
		).
		Record(address).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SelectAddress(addressKey string) (*types.Address, error) {
	sess := c.GetSession()
	address := &types.Address{}

	q := `select * from addresses where address_key = ?;`
	if err := sess.SelectBySql(q).LoadOne(&address); err != nil {
		return nil, fmt.Errorf("failed to select address: %w", err)
	}

	return address, nil
}

func (c *Client) CreatePrice(price *types.Price) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("prices").
		Columns(
			"ticker",
			"price",
		).
		Record(price).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateTransaction(tx *types.Transaction) error {
	sess := c.GetSession()

	if _, err := sess.InsertInto("transactions").
		Columns(
			"user_id",
			"blockchain",
			"tx_hash",
			"tx_key",
			"amount",
			"price",
		).
		Record(tx).Exec(); err != nil {
		return err
	}

	q := `update users set deposits_total = deposits_total + 1, deposits_sum = deposits_sum + ? where id = ?;`
	if _, err := sess.UpdateBySql(q, (tx.Price * tx.Amount), tx.UserID).Exec(); err != nil {
		return err
	}

	return nil
}

func (c *Client) SelectLastPricesByTicker() (map[string]*types.Price, error) {
	sess := c.GetSession()
	res := make([]*types.Price, 0)

	q := `select *
	from prices
	where id in (select max(id)
				 from prices
				 group by ticker);`

	if _, err := sess.SelectBySql(q).Load(&res); err != nil {
		return nil, err
	}

	prices := make(map[string]*types.Price)
	for _, price := range res {
		prices[price.Ticker] = price
	}

	return prices, nil
}
