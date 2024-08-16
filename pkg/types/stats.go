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

	Label          string                    `db:"label" json:"label"`
	Deeplinks      []*ConversionRow          `db:"deeplinks,omitempty" json:"deeplinks,omitempty"`
	DeeplinksLeads map[string]*ConversionRow `db:"-" json:"-"`

	ByDayDB time.Time `db:"by_day" json:"-"`
	ByDay   string    `db:"-" json:"by_day,omitempty"`
}
