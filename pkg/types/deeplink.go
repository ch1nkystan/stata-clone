package types

import (
	"math/rand"
	"time"
)

const (
	DeeplinkLabelReferral = "referral"
)

type Deeplink struct {
	ID    int `db:"id" json:"id"`
	BotID int `db:"bot_id" json:"bot_id"`

	ReferralTelegramID int64  `db:"referral_telegram_id" json:"referral_telegram_id"`
	Hash               string `db:"hash" json:"hash"`
	Label              string `db:"label" json:"label"`

	Active bool `db:"active" json:"active"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

const (
	deeplinkHashLength  = 10
	deeplinkHashSymbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
)

func GenerateDeeplinkHash() string {
	var key string

	rs := []rune(deeplinkHashSymbols)
	lenOfArray := len(rs)

	for i := 0; i < deeplinkHashLength; i++ {
		key += string(rs[rand.Intn(lenOfArray)])
	}
	return key
}
