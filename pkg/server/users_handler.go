package server

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/prosperofair/stata/pkg/types"
)

type UsersSearchRequest struct {
	TelegramID int64  `json:"telegram_id"`
	Username   string `json:"username"`

	ForwardSenderName string `json:"forward_sender_name"`
}

type UsersSearchResponse struct {
	Users []*types.User      `json:"users"`
	Bots  map[int]*types.Bot `json:"bots"`
}

func (req *UsersSearchRequest) validate() error {
	return nil
}

func (s *Server) UsersSearchHandler(c *fiber.Ctx) error {
	req := &UsersSearchRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	res := &UsersSearchResponse{
		Users: make([]*types.User, 0),
		Bots:  make(map[int]*types.Bot),
	}

	// find user by telegram_id
	if req.TelegramID != 0 {
		users, err := s.deps.PG.SelectUsersByTelegramID(req.TelegramID)
		if err != nil {
			return s.InternalServerError(c, err)
		}

		res.Users = users
	} else if req.Username != "" {
		users, err := s.deps.PG.SelectUsersByUsername(req.Username)
		if err != nil {
			return s.InternalServerError(c, err)
		}

		res.Users = users
	} else {
		users, err := s.deps.PG.SelectUsersByForwardSenderName(req.ForwardSenderName)
		if err != nil {
			return s.InternalServerError(c, err)
		}

		res.Users = users
	}

	for _, user := range res.Users {
		if _, ok := res.Bots[user.BotID]; !ok {
			bot, err := s.deps.PG.SelectBotByID(user.BotID)
			if err != nil {
				continue
			}

			res.Bots[user.BotID] = bot
		}
	}

	if len(res.Users) > 5 {
		res.Users = res.Users[:5]
	}

	return c.JSON(res)
}

type UsersGetRequest struct {
	BotToken   string `json:"bot_token"`
	TelegramID int64  `json:"telegram_id"`
}

type UsersGetResponse struct {
	User *types.User `json:"user"`
}

func (req *UsersGetRequest) validate() error {
	return nil
}

func (s *Server) UsersGetHandler(c *fiber.Ctx) error {
	req := &UsersGetRequest{}
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

	users, err := s.deps.PG.SelectBotUsersByTelegramID(bot.ID, req.TelegramID)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	if len(users) == 0 {
		return s.BadRequest(c, errors.New("user not found"))
	}

	res := &UsersGetResponse{
		User: users[0],
	}

	return c.JSON(res)
}

type UsersSetDefaultChannelRequest struct {
	BotToken string `json:"bot_token"`

	DepotChannelHash   string `json:"depot_channel_hash"`
	TelegramChannelID  int64  `json:"telegram_channel_id"`
	TelegramChannelURL string `json:"telegram_channel_url"`
}

func (req *UsersSetDefaultChannelRequest) validate() error {
	if req.TelegramChannelID == 0 {
		return errors.New("telegram_channel_id is empty")
	}

	if req.TelegramChannelURL == "" {
		return errors.New("telegram_channel_url is empty")
	}

	if req.DepotChannelHash == "" {
		return errors.New("depot_channel_hash is empty")
	}

	return nil
}

func (s *Server) UsersSetDefaultChannelHandler(c *fiber.Ctx) error {
	req := &UsersSetDefaultChannelRequest{}
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

	if err := s.deps.PG.UpdateBotUsersSetDefaultChannel(
		bot.ID,
		req.DepotChannelHash,
		req.TelegramChannelID,
		req.TelegramChannelURL,
	); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}

type UsersUpdateTelegramChannelRequest struct {
	BotToken         string `json:"bot_token"`
	DepotChannelHash string `json:"depot_channel_hash"`

	NewTelegramChannelID  int64  `json:"new_telegram_channel_id"`
	NewTelegramChannelURL string `json:"new_telegram_channel_url"`
}

func (req *UsersUpdateTelegramChannelRequest) validate() error {
	if req.BotToken == "" {
		return errors.New("bot_token is empty")
	}

	if req.NewTelegramChannelID == 0 {
		return errors.New("new_telegram_channel_id is empty")
	}

	if req.NewTelegramChannelURL == "" {
		return errors.New("new_telegram_channel_url is empty")
	}

	return nil
}

func (s *Server) UsersUpdateTelegramChannelHandler(c *fiber.Ctx) error {
	req := &UsersUpdateTelegramChannelRequest{}
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

	if err := s.deps.PG.UpdateBotUsersTelegramChannel(bot.ID, req.DepotChannelHash,
		req.NewTelegramChannelID, req.NewTelegramChannelURL,
	); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}
