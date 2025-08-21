package server

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/stata/pkg/types"
)

type reportsFilterRequest struct {
	BotToken string    `json:"bot_token"`
	StartAt  time.Time `json:"start_at"`
	EndAt    time.Time `json:"end_at"`
}

func (req *reportsFilterRequest) validate() error {
	if req.BotToken == "" {
		return errors.New("bot_token is empty")
	}

	return nil
}

type reportsLeadsByCampaignResponse struct {
	Data []*types.LeadsByCampaignRow `json:"data"`
}

func (s *Server) reportsLeadsByCampaign(c *fiber.Ctx) error {
	req := &reportsFilterRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	data, err := s.deps.PG.SelectLeadsByCampaign(req.BotToken, req.StartAt, req.EndAt)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res := &reportsLeadsByCampaignResponse{
		Data: data,
	}

	return c.JSON(res)
}
