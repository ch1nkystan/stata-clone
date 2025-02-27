package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type onlineSnapshotResponse struct {
	Snapshot map[time.Time]int `json:"snapshot"`
}

func (s *Server) onlineSnapshotHandler(c *fiber.Ctx) error {
	req := &dateRangeRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	snapshot, err := s.deps.PG.SelectOnlineSnapshotForInterval(bot.ID, req.Start, req.End)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res := &onlineSnapshotResponse{
		Snapshot: snapshot,
	}

	return c.JSON(res)
}
