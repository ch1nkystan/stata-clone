package types

import "time"

type ConversionRow struct {
	UsersTotal      int     `db:"users_total" json:"users_total"`
	UsersUnique     int     `db:"users_unique" json:"users_unique"`
	UsersUniqueRate float64 `db:"users_unique_rate" json:"users_unique_rate"`

	LeadsTotal          int     `db:"leads_total" json:"leads_total"`
	LeadsUsers          int     `db:"leads_users" json:"leads_users"`
	LeadsPerUser        float64 `db:"leads_per_user" json:"leads_per_user"`
	LeadsConversionRate float64 `db:"leads_conversion_rate" json:"leads_conversion_rate"`
	Profit              float64 `db:"profit" json:"profit"`
	Income              float64 `db:"income" json:"income"`
	Expenses            float64 `db:"expenses" json:"expenses"`

	Label          string                    `db:"label" json:"label"`
	Deeplinks      []*ConversionRow          `db:"deeplinks,omitempty" json:"deeplinks,omitempty"`
	DeeplinksLeads map[string]*ConversionRow `db:"-" json:"-"`

	ByDayDB time.Time `db:"by_day" json:"-"`
	ByDay   string    `db:"-" json:"by_day,omitempty"`
}

type DepositRow struct {
	ID         int       `db:"id" json:"id"`
	Hash       string    `db:"hash" json:"hash"`
	Deeplink   string    `db:"deeplink" json:"deeplink"`
	Blockchain string    `db:"blockchain" json:"blockchain"`
	Amount     float64   `db:"amount" json:"amount"`
	Date       time.Time `db:"date" json:"date"`
}

type UsersRecap struct {
	UsersTotal  int     `db:"users_total" json:"users_total"`
	UsersUnique int     `db:"users_unique" json:"users_unique"`
	UniqueRate  float64 `db:"unique_rate" json:"unique_rate"`
}
