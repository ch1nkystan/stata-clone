package types

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
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

type PixelLink struct {
	ID int `json:"id" db:"id"`

	FBAccessMarker string `json:"fb_access_marker" db:"fb_access_marker"`
	FBPixelID      int64  `json:"fb_pixel_id" db:"fb_pixel_id"`
	FBC            string `json:"fbc" db:"fbc"`
	FBP            string `json:"fbp" db:"fbp"`
	DeeplinkID     int    `json:"deeplink_id" db:"deeplink_id"`
	DeeplinkHash   string `json:"deeplink_hash" db:"deeplink_hash"`

	InviteUUID uuid.UUID `json:"invite_uuid" db:"invite_uuid"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
