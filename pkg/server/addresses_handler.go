package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/stata/pkg/types"
)

type addressesCreateRequest struct {
	Blockchain string `json:"blockchain"`
	AddressKey string `json:"address_key"`
	Address    string `json:"address"`
	BID        string `json:"bid"`
}

func (req *addressesCreateRequest) validate() error {
	whitelist := map[string]struct{}{
		types.BlockchainBSC: {},
		types.BlockchainETH: {},
		types.BlockchainBTC: {},
	}

	if _, ok := whitelist[req.Blockchain]; !ok {
		return errors.New("invalid blockchain")
	}

	return nil
}

func (s *Server) addressesCreateHandler(c *fiber.Ctx) error {
	req := &addressesCreateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	address := &types.Address{
		Blockchain: req.Blockchain,
		AddressKey: req.AddressKey,
		Address:    req.Address,
		BID:        req.BID,
	}

	if err := s.deps.PG.CreateAddress(address); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}

type addressesCheckRequest struct {
	AddressKey string `json:"address_key"`
}

func (req *addressesCheckRequest) validate() error {
	if req.AddressKey == "" {
		return errors.New("address_key is required")
	}

	return nil
}

func (s *Server) addressesCheckHandler(c *fiber.Ctx) error {
	req := &addressesCheckRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	if _, err := s.deps.PG.SelectAddress(req.AddressKey); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}
