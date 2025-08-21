package types

const (
	EventTypeRegister = "register"
	EventTypeMessage  = "message"
	EventTypeLaunch   = "launch"
	EventTypeMailing  = "mailing"
	EventTypeDeposit  = "deposit"
)

var EventsWhitelist = map[string]struct{}{
	EventTypeRegister: {},
	EventTypeMessage:  {},
	EventTypeMailing:  {},
	EventTypeDeposit:  {},
}

var Events = []string{EventTypeRegister, EventTypeMessage, EventTypeMailing, EventTypeDeposit}

type Event struct {
	ID        int    `db:"id" json:"id"`
	EventType string `db:"event_type" json:"event_type"`

	TelegramID int64  `db:"telegram_id" json:"telegram_id"`
	Firstname  string `db:"firstname" json:"firstname"`
	Lastname   string `db:"lastname" json:"lastname"`
	Username   string `db:"username" json:"username"`

	IsBot        bool   `db:"is_bot" json:"is_bot"`
	IsPremium    bool   `db:"is_premium" json:"is_premium"`
	LanguageCode string `db:"language_code" json:"language_code"`
}

type EventsLog struct {
	ID        int    `db:"id" json:"id"`
	EventType string `db:"event_type" json:"event_type"`

	ReporterTelegramID int64 `db:"reporter_telegram_id" json:"reporter_telegram_id"`
	UserID             int   `db:"user_id" json:"user_id"`
	TelegramID         int64 `db:"telegram_id" json:"telegram_id"`

	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}
