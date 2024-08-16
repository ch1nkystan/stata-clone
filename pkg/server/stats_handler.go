package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/types"
	"go.uber.org/zap"
)

type StatsUsersCount struct {
	Total      int `json:"total"`
	Success    int `json:"success"`
	Indefinite int `json:"indefinite"`
	Fail       int `json:"fail"`
}

type StatsMailingStateRequest struct {
	BotToken string `json:"bot_token"`
}

type StatsMailingStateResponse struct {
	Channels map[string]StatsUsersCount `json:"channels"`
}

func (s *Server) StatsMailingStateHandler(c *fiber.Ctx) error {
	req := &StatsMailingStateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	stats, err := s.deps.PG.SelectBotMailingStats(bot.ID)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res := StatsMailingStateResponse{
		Channels: make(map[string]StatsUsersCount),
	}

	for _, stat := range stats {
		log.Info("stat", zap.Any("stat", stat))
		current, ok := res.Channels[stat.DepotChannelHash]
		if !ok {
			res.Channels[stat.DepotChannelHash] = StatsUsersCount{}
		}

		current.Total += stat.Total
		current.Success += stat.Success
		current.Indefinite += stat.Indefinite
		current.Fail += stat.Fail

		res.Channels[stat.DepotChannelHash] = current
	}

	return c.JSON(res)
}

type StatsUsersCountRequest struct {
}

type StatsUsersCountResponse struct {
	Bots map[string]StatsUsersCount `json:"bots"`
}

func (s *Server) StatsUsersCountHandler(c *fiber.Ctx) error {
	req := &StatsUsersCountRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	bots, err := s.deps.PG.SelectAllBots()
	if err != nil {
		return s.InternalServerError(c, err)
	}

	botsByID := make(map[int]*types.Bot)
	for _, bot := range bots {
		botsByID[bot.ID] = bot
	}

	stats, err := s.deps.PG.SelectUsersCountStats()
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res := StatsUsersCountResponse{
		Bots: make(map[string]StatsUsersCount),
	}

	for _, stat := range stats {
		bot, ok := botsByID[stat.BotID]
		if !ok {
			continue
		}

		res.Bots[bot.BotToken] = StatsUsersCount{
			Total:      stat.Total,
			Success:    stat.Success,
			Indefinite: stat.Indefinite,
			Fail:       stat.Fail,
		}
	}

	return c.JSON(res)
}
