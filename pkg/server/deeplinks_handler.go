package server

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/stata/pkg/types"
)

type DeeplinksCreateRequest struct {
	BotToken           string `json:"bot_token"`
	ReferralTelegramID int64  `json:"referral_telegram_id"`
	Label              string `json:"label"`
	Hash               string `json:"hash,omitempty"`
}

type DeeplinksCreateResponse struct {
	Hash string `json:"hash"`
}

func (req *DeeplinksCreateRequest) validate() error {
	if req.Label == "" {
		return errors.New("invalid label")
	}

	return nil
}

func (s *Server) DeeplinksCreateHandler(c *fiber.Ctx) error {
	req := &DeeplinksCreateRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c,
			fmt.Errorf("select bot by token: %w", err))
	}

	res := &DeeplinksCreateResponse{}
	deeplink := &types.Deeplink{
		BotID:              bot.ID,
		ReferralTelegramID: req.ReferralTelegramID,
		Hash:               types.GenerateDeeplinkHash(),
		Label:              req.Label,
	}

	if req.Hash != "" {
		deeplinks, err := s.deps.PG.SelectBotDeeplinksByHash(bot.ID, req.Hash)
		if err != nil {
			return s.InternalServerError(c, err)
		}

		if len(deeplinks) > 0 {
			res.Hash = deeplinks[0].Hash
			return c.JSON(res)
		}

		deeplink.Hash = req.Hash
	}

	if deeplink.ReferralTelegramID > 0 {
		deeplink.Label = types.DeeplinkLabelReferral
		deeplinks, err := s.deps.PG.SelectBotDeeplinksByReferralID(bot.ID, req.ReferralTelegramID)
		if err != nil {
			return s.InternalServerError(c, err)
		}

		if len(deeplinks) > 0 {
			res.Hash = deeplinks[0].Hash
			return c.JSON(res)
		}
	}

	if err := s.deps.PG.CreateDeeplink(deeplink); err != nil {
		return s.InternalServerError(c, err)
	}

	return c.JSON(DeeplinksCreateResponse{
		Hash: deeplink.Hash,
	})
}

type DeeplinksListRequest struct {
	BotToken string `json:"bot_token"`
}

type DeeplinksListResponse struct {
	Deeplinks []*types.Deeplink `json:"deeplinks"`
}

func (s *Server) DeeplinksListHandler(c *fiber.Ctx) error {
	req := &DeeplinksListRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.InternalServerError(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c,
			fmt.Errorf("select bot by token: %w", err))
	}

	deeplinks, err := s.deps.PG.SelectBotDeeplinksByReferralID(bot.ID, 0)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	return c.JSON(DeeplinksListResponse{
		Deeplinks: deeplinks,
	})
}

type DeeplinksUpdateRequest struct {
	BotToken string `json:"bot_token"`
	Hash     string `json:"hash"`
	Label    string `json:"label"`
}

func (req *DeeplinksUpdateRequest) validate() error {
	if req.Label == "" {
		return errors.New("invalid label")
	}

	return nil
}

func (s *Server) DeeplinksUpdateHandler(c *fiber.Ctx) error {
	req := &DeeplinksUpdateRequest{}
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

	if err := s.deps.PG.UpdateDeeplinkLabel(bot.ID, req.Hash, req.Label); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}
