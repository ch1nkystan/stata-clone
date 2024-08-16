package main

import (
	"fmt"
	"time"

	"github.com/gocraft/dbr/v2"
)

type IncomeInfo struct {
	UserID       int64     `db:"user_id"`
	BotLink      string    `db:"bot_link"`
	BotName      string    `db:"bot_name"`
	IncomeSource string    `db:"income_source"`
	TypeBot      string    `db:"type_bot"`
	CreateAt     time.Time `db:"create_at"`
}

type PGSQLClient struct {
	pgdb *dbr.Connection
}

func NewPGSQLClient(conn *dbr.Connection) *PGSQLClient {
	return &PGSQLClient{
		pgdb: conn,
	}
}

func (c *PGSQLClient) GetSession() *dbr.Session {
	return c.pgdb.NewSession(nil)
}

func (c *PGSQLClient) SelectIncomesInfo(limit int) ([]*IncomeInfo, error) {
	sess := c.GetSession()

	res := make([]*IncomeInfo, 0)

	q := `select * from income_info order by create_at desc limit ?`
	if _, err := sess.SelectBySql(q, limit).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *PGSQLClient) SelectIncomesInfoForUser(telegramID int64, botUsername string) ([]*IncomeInfo, error) {
	sess := c.GetSession()

	res := make([]*IncomeInfo, 0)

	botLink := fmt.Sprintf("https://t.me/%s", botUsername)

	q := `select * from income_info where user_id = ? and bot_link = ? order by create_at asc`
	if _, err := sess.SelectBySql(q, telegramID, botLink).Load(&res); err != nil {
		return nil, err
	}

	return res, nil
}
