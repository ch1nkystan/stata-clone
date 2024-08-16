package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/stata/pkg/types"
)

type TransactionsCreateRequest struct {
	UserID     int     `json:"user_id"`
	Amount     float64 `json:"amount"`
	Blockchain string  `json:"blockchain"`
	TXHash     string  `json:"tx_hash"`
	TXKey      string  `json:"tx_key"`
}

func (req *TransactionsCreateRequest) validate() error {
	whitelist := map[string]struct{}{
		types.TransactionBlockchainBSC: {},
		types.TransactionBlockchainETH: {},
		types.TransactionBlockchainBTC: {},
	}

	if _, ok := whitelist[req.Blockchain]; !ok {
		return errors.New("invalid blockchain")
	}

	return nil
}

func (s *Server) TransactionsCreateHandler(c *fiber.Ctx) error {
	req := &TransactionsCreateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	tx := &types.Transaction{
		Amount:     req.Amount,
		Blockchain: req.Blockchain,
		TXHash:     req.TXHash,
		TXKey:      req.TXKey,
	}

	user, err := s.deps.PG.SelectUserByID(req.UserID)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	tx.UserID = user.ID

	prices, err := s.deps.PG.SelectLastPricesByTicker()
	if err != nil {
		return s.InternalServerError(c, err)
	}

	tickerMap := map[string]string{
		types.TransactionBlockchainBSC: "BNBUSDT",
		types.TransactionBlockchainETH: "ETHUSDT",
		types.TransactionBlockchainBTC: "BTCUSDT",
	}

	if price, ok := prices[tickerMap[tx.Blockchain]]; ok {
		tx.Price = price.Price
	}

	if err := s.deps.PG.CreateTransaction(tx); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}
