package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"github.com/prosperofair/stata/pkg/types"
	"go.uber.org/zap"
)

type Ticker struct {
	HTTPClient *http.Client
	pg         *pgsql.Client

	tick   time.Duration
	aid    int
	symbol string

	dryRun bool
}

func NewTicker(pg *pgsql.Client, symbol string, tick time.Duration, dryRun bool) *Ticker {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return &Ticker{
		HTTPClient: client,
		pg:         pg,

		tick:   tick,
		symbol: symbol,

		dryRun: dryRun,
	}
}

type BinanceResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"-"`
}

func (m *BinanceResponse) UnmarshalJSON(data []byte) error {
	type Alias BinanceResponse
	aux := &struct {
		*Alias
		FloatValStr string `json:"price"`
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if f, err := strconv.ParseFloat(aux.FloatValStr, 64); err == nil {
		m.Price = f
	} else {
		return err
	}
	return nil
}

func (t *Ticker) run() error {
	u := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", t.symbol)
	for {
		func() {
			resp, err := t.HTTPClient.Get(u)
			if err != nil {
				log.Error("failed to get btc price", zap.Error(err))
				return
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error("failed to read body", zap.Error(err))
				return
			}

			response := BinanceResponse{}
			if err := json.Unmarshal(body, &response); err != nil {
				log.Error("failed to unmarshal body",
					zap.Error(err),
					zap.String("body", string(body)),
					zap.String("symbol", t.symbol),
				)
				return
			}

			if !t.dryRun {
				if err := t.pg.CreatePrice(&types.Price{
					Ticker: t.symbol,
					Price:  response.Price,
				}); err != nil {
					log.Error("failed to insert btc price", zap.Error(err))
					return
				}
			}

			log.Info("tick", zap.String("symbol", response.Symbol),
				zap.Float64("price", response.Price))
		}()

		time.Sleep(t.tick)
	}
}
