package server

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/stata/pkg/types"
)

type reportsLeadsByCampaignRequest struct {
	BotToken string `json:"bot_token"`

	StartAt string    `json:"start_at"`
	Start   time.Time `json:"-"`

	EndAt string    `json:"end_at"`
	End   time.Time `json:"-"`
}

func (req *reportsLeadsByCampaignRequest) validate() error {
	if req.BotToken == "" {
		return errors.New("bot_token is empty")
	}

	sat, err := time.Parse(time.DateOnly, req.StartAt)
	if err != nil {
		sat = time.Now().AddDate(0, -1, 0)
	}
	req.Start = sat

	eat, err := time.Parse(time.DateOnly, req.EndAt)
	if err != nil {
		eat = time.Now()
	}
	req.End = eat

	return nil
}

type reportsLeadsByCampaignResponse struct {
	Leads []*types.LeadByCampaignRow `json:"leads"`
}

func (s *Server) reportsLeadsByCampaignHandler(c *fiber.Ctx) error {
	req := &reportsLeadsByCampaignRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	leads, err := s.deps.PG.SelectLeadsByCampaign(req.BotToken, req.Start, req.End)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res := &reportsLeadsByCampaignResponse{
		Leads: leads,
	}

	return c.JSON(res)
}
