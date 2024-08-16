package types

import "time"

type Bot struct {
	ID int `db:"id" json:"id"`

	APIKey      string `db:"api_key" json:"api_key"`
	BotToken    string `db:"bot_token" json:"bot_token"`
	BotUsername string `db:"bot_username" json:"bot_username"`
	BotType     string `db:"bot_type" json:"bot_type"`
	BID         string `db:"bid" json:"bid"`

	Active bool `db:"active" json:"active"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
