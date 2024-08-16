package server

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/types"
	"go.uber.org/zap"
)

type BotsRegisterRequest struct {
	APIKey      string `json:"api_key"`
	BotUsername string `json:"bot_username"`
	BotType     string `json:"bot_type"`
	BuyerID     string `json:"bid"`
}

type BotsRegisterResponse struct {
}

func (req *BotsRegisterRequest) validate() error {
	return nil
}

func (s *Server) BotsRegisterHandler(c *fiber.Ctx) error {
	req := &BotsRegisterRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	hash := md5.Sum([]byte(req.APIKey))

	bot := &types.Bot{
		APIKey:      req.APIKey,
		BotToken:    hex.EncodeToString(hash[:]),
		BotUsername: req.BotUsername,
		BotType:     req.BotType,
		BID:         req.BuyerID,
	}

	if err := s.deps.PG.CreateBot(bot); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}

type BotsImportRequest struct {
	Bots []BotsRegisterRequest `json:"bots"`
}

type BotsImportResponse struct {
	Affected int `json:"affected"`
}

func (s *Server) BotsImportHandler(c *fiber.Ctx) error {
	req := &BotsImportRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	res := &BotsImportResponse{
		Affected: 0,
	}

	for _, bot := range req.Bots {
		if err := bot.validate(); err != nil {
			return s.BadRequest(c, err)
		}

		hash := md5.Sum([]byte(bot.APIKey))

		if err := s.deps.PG.CreateBot(&types.Bot{
			APIKey:      bot.APIKey,
			BotToken:    hex.EncodeToString(hash[:]),
			BotUsername: bot.BotUsername,
			BotType:     bot.BotType,
		}); err != nil {
			log.Error("failed to create bot", zap.Error(err))
			continue
		}

		res.Affected++
	}

	return c.JSON(res)
}
