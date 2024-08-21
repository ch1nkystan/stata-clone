package types

import "time"

type FBToolToken struct {
	ID    int    `json:"id" db:"id"`
	Token string `json:"token" db:"token"`

	Active      bool `json:"active" db:"active"`
	DaysToFetch int  `json:"days_to_fetch" db:"days_to_fetch"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	FetchedAt time.Time `json:"fetched_at" db:"fetched_at"`
}

type FBToolAccount struct {
	ID int `json:"id" db:"id"`

	TokenID           int    `json:"token_id" db:"token_id"`
	FBToolAccountID   int    `json:"fbtool_account_id" db:"fbtool_account_id"`
	FBToolAccountName string `json:"fbtool_account_name" db:"fbtool_account_name"`

	Active bool `json:"active" db:"active"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	FetchedAt time.Time `json:"fetched_at" db:"fetched_at"`
}

type FBToolCampaignStat struct {
	ID int `json:"id" db:"id"`

	FBToolAccountID int `json:"fbtool_account_id" db:"fbtool_account_id"`

	CampaignName string `json:"campaign_name" db:"campaign_name"`
	CampaignID   string `json:"campaign_id" db:"campaign_id"`

	Status          string `json:"status" db:"status"`
	EffectiveStatus string `json:"effective_status" db:"effective_status"`

	Impressions int       `json:"impressions" db:"impressions"`
	Clicks      int       `json:"clicks" db:"clicks"`
	Spend       float64   `json:"spend" db:"spend"`
	Date        time.Time `json:"date" db:"date"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
