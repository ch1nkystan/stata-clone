package types

import "time"

const (
	BlockchainBSC = "bsc"
	BlockchainETH = "eth"
	BlockchainBTC = "btc"
)

type Address struct {
	ID int `db:"id" json:"id"`

	Blockchain string `db:"blockchain" json:"blockchain"`
	AddressKey string `db:"address_key" json:"address_key"`
	Address    string `db:"address" json:"address"`

	BID string `db:"bid" json:"bid"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Transaction struct {
	ID int `db:"id" json:"id"`

	UserID int `db:"user_id" json:"user_id"`

	Blockchain string `db:"blockchain" json:"blockchain"`
	TXHash     string `db:"tx_hash" json:"tx_hash"`
	TXKey      string `db:"tx_key" json:"tx_key"`

	Amount float64 `db:"amount" json:"amount"`
	Price  float64 `db:"price" json:"price"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Price struct {
	ID        int       `db:"id" json:"id"`
	Ticker    string    `db:"ticker" json:"ticker"`
	Price     float64   `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Range struct {
	Start time.Time
	End   time.Time
}
