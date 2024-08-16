package types

import "time"

const (
	UserMailingMaxFailedAttempts = 3

	UserMailingStateReady      = "ready"
	UserMailingStateInProgress = "in-progress"
	UserMailingStateFinished   = "finished"
	UserMailingStateBlocked    = "blocked"
)

type User struct {
	ID         int `db:"id" json:"id"`
	BotID      int `db:"bot_id" json:"bot_id"`
	DeeplinkID int `db:"deeplink_id" json:"deeplink_id"`

	DepotChannelHash   string `db:"depot_channel_hash" json:"depot_channel_hash"`
	TelegramChannelID  int64  `db:"telegram_channel_id" json:"telegram_channel_id"`
	TelegramChannelURL string `db:"telegram_channel_url" json:"telegram_channel_url"`

	TelegramID        int64  `db:"telegram_id" json:"telegram_id"`
	Firstname         string `db:"first_name" json:"first_name"`
	Lastname          string `db:"last_name" json:"last_name"`
	ForwardSenderName string `db:"forward_sender_name" json:"forward_sender_name"`
	Username          string `db:"username" json:"username"`

	IsBot        bool   `db:"is_bot" json:"is_bot"`
	IsPremium    bool   `db:"is_premium" json:"is_premium"`
	LanguageCode string `db:"language_code" json:"language_code"`

	DepositsTotal int       `db:"deposits_total" json:"deposits_total"`
	DepositsSum   float64   `db:"deposits_sum" json:"deposits_sum"`
	Deposited     bool      `db:"deposited" json:"deposited"`
	DepositedAt   time.Time `db:"deposited_at" json:"deposited_at"`

	Seen         int    `db:"seen" json:"seen"`
	Active       bool   `db:"active" json:"active"`
	EventCreated string `db:"event_created" json:"event_created"`

	MailingState          string    `db:"mailing_state" json:"mailing_state"`
	MailingStateUpdatedAt time.Time `db:"mailing_state_updated_at" json:"mailing_state_updated_at"`
	MailingFailedAttempts int       `db:"mailing_failed_attempts" json:"mailing_failed_attempts"`

	MessagedAt time.Time `db:"messaged_at" json:"messaged_at"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

func (u *User) GenerateForwardSenderName() string {
	res := u.Firstname
	if u.Lastname != "" {
		res += " " + u.Lastname
	}

	return res
}
