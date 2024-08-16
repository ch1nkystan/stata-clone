package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/prosperofair/stata/pkg/types"
)

type MailingPrepareUsersListRequest struct {
	BotToken         string `json:"bot_token"`
	DepotChannelHash string `json:"depot_channel_hash"`

	Limit  int  `json:"limit"`
	DryRun bool `json:"dry_run"`
}

type MailingPrepareUsersListResponse struct {
	Users []*types.User `json:"users"`
}

func (req *MailingPrepareUsersListRequest) validate() error {
	if req.BotToken == "" {
		return errors.New("bot_token is empty")
	}

	if req.DepotChannelHash == "" {
		return errors.New("depot_channel_hash is empty")
	}

	if req.Limit == 0 {
		return errors.New("limit is empty")
	}

	return nil
}

func (s *Server) MailingPrepareUsersListHandler(c *fiber.Ctx) error {
	req := &MailingPrepareUsersListRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	users, err := s.deps.PG.SelectRandomReadyUsersByDepotChannelHash(bot.ID, req.DepotChannelHash, req.Limit)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res := &MailingPrepareUsersListResponse{
		Users: users,
	}

	if req.DryRun {
		return c.JSON(res)
	}

	ids := make([]int, 0)
	for _, u := range users {
		ids = append(ids, u.ID)
	}
	if err := s.deps.PG.UpdateUsersMailingState(types.UserMailingStateInProgress, ids); err != nil {
		return s.InternalServerError(c, err)
	}

	return c.JSON(res)
}

type MailingUpdateUserStateRequest struct {
	BotToken   string `json:"bot_token"`
	TelegramID int64  `json:"telegram_id"`
	State      string `json:"state"`
}

func (req *MailingUpdateUserStateRequest) validate() error {
	if req.BotToken == "" {
		return errors.New("bot_token is empty")
	}

	if req.TelegramID == 0 {
		return errors.New("telegram_id is empty")
	}

	whitelist := map[string]struct{}{
		types.UserMailingStateReady:      {},
		types.UserMailingStateInProgress: {},
		types.UserMailingStateFinished:   {},
		types.UserMailingStateBlocked:    {},
	}

	if _, ok := whitelist[req.State]; !ok {
		return errors.New("invalid state")
	}

	return nil
}

func (s *Server) MailingUpdateUserStateHandler(c *fiber.Ctx) error {
	req := &MailingUpdateUserStateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	if err := s.deps.PG.UpdateBotUserMailingState(req.BotToken, req.TelegramID, req.State); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}

type MailingFinishUsersListRequest struct {
	BotToken         string `json:"bot_token"`
	DepotChannelHash string `json:"depot_channel_hash"`
}

func (req *MailingFinishUsersListRequest) validate() error {
	if req.BotToken == "" {
		return errors.New("bot_token is empty")
	}

	if req.DepotChannelHash == "" {
		return errors.New("depot_channel_hash is empty")
	}

	return nil
}

func (s *Server) MailingFinishUsersListHandler(c *fiber.Ctx) error {
	req := &MailingFinishUsersListRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	if err := s.deps.PG.SetBotUsersMailingStatesReady(bot.ID, req.DepotChannelHash); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}
