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
	Expense             float64 `db:"expense" json:"expense"`
	Clicks              int     `db:"clicks" json:"clicks"`
	Impressions         int     `db:"impressions" json:"impressions"`

	Label string `db:"label" json:"label,omitempty"`

	ByDayDB time.Time `db:"by_day" json:"-"`
	ByDay   string    `db:"-" json:"by_day,omitempty"`

	ByPeriodDB  time.Time `db:"by_period" json:"-"`
	PeriodStart string    `db:"-" json:"period_start,omitempty"`
	PeriodEnd   string    `db:"-" json:"period_end,omitempty"`
}

type DepositRow struct {
	ID         int       `db:"id" json:"id"`
	Hash       string    `db:"hash" json:"hash"`
	Deeplink   string    `db:"deeplink" json:"deeplink"`
	Blockchain string    `db:"blockchain" json:"blockchain"`
	Amount     float64   `db:"amount" json:"amount"`
	Date       time.Time `db:"date" json:"date"`
}

type MetricRow struct {
	AllTime    interface{} `db:"all_time" json:"all_time"`
	Period     interface{} `db:"period" json:"period"`
	LastPeriod interface{} `db:"last_period" json:"last_period"`
	Diff       interface{} `db:"diff" json:"diff"`
}

type Snapshot struct {
	ID        int       `db:"id" json:"id"`
	Snapshot  string    `db:"snapshot" json:"snapshot"`
	BotID     int       `db:"bot_id" json:"bot_id"`
	Users     int       `db:"users" json:"users"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type LeadByCampaignRow struct {
	Label             string  `db:"label" json:"label"`
	Total             int     `db:"total" json:"total"`
	Unique            int     `db:"unique" json:"unique"`
	UniquenessRatio   float64 `db:"uniquenes_ratio" json:"uniquenes_ratio"`
	Deposited         int     `db:"deposited" json:"deposited"`
	DepositedRatio    float64 `db:"deposited_ratio" json:"deposited_ratio"`
	TotalSpend        float64 `db:"total_spend" json:"total_spend"`
	DepositsSum       float64 `db:"deposits_sum" json:"deposits_sum"`
	DepositsTotal     int     `db:"deposits_total" json:"deposits_total"`
	CPU               float64 `db:"cpu" json:"cpu"`
	DepositsPerPerson float64 `db:"deposits_per_person" json:"deposits_per_person"`
	DepositPrice      float64 `db:"deposit_price" json:"deposit_price"`
	PNL               float64 `db:"pnl" json:"pnl"`
}
