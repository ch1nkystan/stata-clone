package server

import (
	"errors"
	"fmt"
	"net"

	"github.com/avct/uasurfer"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/types"
)

type EventsSubmitUserRegisterRequest struct {
	BotToken string `json:"bot_token"`

	TelegramID int64  `json:"telegram_id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Username   string `json:"username"`

	IsBot        bool   `json:"is_bot"`
	IsPremium    bool   `json:"is_premium"`
	LanguageCode string `json:"language_code"`

	Hash string `json:"hash"`
}

func (req *EventsSubmitUserRegisterRequest) validate() error {
	if req.TelegramID == 0 {
		return errors.New("invalid telegram_id")
	}

	return nil
}

func (s *Server) EventsSubmitUserRegisterHandler(c *fiber.Ctx) error {
	req := &EventsSubmitUserRegisterRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c,
			fmt.Errorf("select bot by token: %w", err))
	}

	users, err := s.deps.PG.SelectBotUsersByTelegramID(bot.ID, req.TelegramID)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	metricsEvents.WithLabelValues(bot.BotToken, types.EventTypeRegister).Inc()

	deeplinkID := 0
	deeplinks, err := s.deps.PG.SelectBotDeeplinksByHash(bot.ID, req.Hash)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	if len(deeplinks) > 0 {
		deeplinkID = deeplinks[0].ID
	}

	if len(users) > 0 {
		user := users[0]
		if deeplinkID != 0 && users[0].DeeplinkID == 0 {
			if err := s.deps.PG.UpdateUserDeeplink(user.ID, deeplinkID); err != nil {
				return s.InternalServerError(c, err)
			}
		}
	} else {
		seenTimes, err := s.deps.PG.CountUsersByTelegramID(req.TelegramID)
		if err != nil {
			return s.InternalServerError(c, err)
		}

		user := &types.User{
			BotID:        bot.ID,
			DeeplinkID:   deeplinkID,
			TelegramID:   req.TelegramID,
			Firstname:    req.Firstname,
			Lastname:     req.Lastname,
			Username:     req.Username,
			IsBot:        req.IsBot,
			IsPremium:    req.IsPremium,
			LanguageCode: req.LanguageCode,
			EventCreated: types.EventTypeRegister,
			Seen:         seenTimes,
		}

		user.ForwardSenderName = user.GenerateForwardSenderName()

		channel, err := s.deps.Depot.GetBotChannelByTDR(bot.BotToken)
		if err != nil {
			log.Error("failed to list bot active channels", zap.Error(err))
		}

		if channel != nil {
			user.DepotChannelHash = channel.Hash
			user.TelegramChannelID = channel.TelegramChannelID
			user.TelegramChannelURL = channel.TelegramChannelURL
		}

		log.Info("creating user",
			zap.String("bot_token", bot.BotToken),
			zap.Any("channel", channel),
			zap.Int64("telegram_id", user.TelegramID),
		)

		if err := s.deps.PG.CreateUser(user); err != nil {
			return s.InternalServerError(c, err)
		}
	}

	return s.ResponseOK(c)
}

type EventsSubmitMessageRequest struct {
	BotToken string `json:"bot_token"`

	TelegramID int64  `json:"telegram_id"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Username   string `json:"username"`

	IsBot        bool   `json:"is_bot"`
	IsPremium    bool   `json:"is_premium"`
	LanguageCode string `json:"language_code"`

	Subscribed bool `json:"subscribed"`
}

func (req *EventsSubmitMessageRequest) validate() error {
	if req.TelegramID == 0 {
		return errors.New("invalid telegram_id")
	}

	return nil
}

func (s *Server) EventsSubmitMessageHandler(c *fiber.Ctx) error {
	req := &EventsSubmitMessageRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c,
			fmt.Errorf("select bot by token: %w", err))
	}

	metricsEvents.WithLabelValues(bot.BotToken, types.EventTypeMessage).Inc()

	users, err := s.deps.PG.SelectBotUsersByTelegramID(bot.ID, req.TelegramID)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	user := &types.User{
		BotID:        bot.ID,
		TelegramID:   req.TelegramID,
		Firstname:    req.Firstname,
		Lastname:     req.Lastname,
		Username:     req.Username,
		IsBot:        req.IsBot,
		IsPremium:    req.IsPremium,
		LanguageCode: req.LanguageCode,
		EventCreated: types.EventTypeMessage,
		Subscribed:   req.Subscribed,
	}

	user.ForwardSenderName = user.GenerateForwardSenderName()

	if len(users) > 0 {
		oldUser := users[0]
		if oldUser.DepotChannelHash == "" {
			channel, err := s.deps.Depot.GetBotChannelByTDR(bot.BotToken)
			if err != nil {
				log.Error("failed to list bot active channels", zap.Error(err))
			}

			if channel != nil {
				user.DepotChannelHash = channel.Hash
				user.TelegramChannelID = channel.TelegramChannelID
				user.TelegramChannelURL = channel.TelegramChannelURL
			}
		}

		if err := s.deps.PG.UpdateUserOnMessage(users[0], user); err != nil {
			return s.InternalServerError(c, err)
		}

		log.Info("user updated", zap.Time("messaged_at", users[0].MessagedAt), zap.Int("id", users[0].ID))
	} else {
		channel, err := s.deps.Depot.GetBotChannelByTDR(bot.BotToken)
		if err != nil {
			log.Error("failed to list bot active channels", zap.Error(err))
		}

		if channel != nil {
			user.DepotChannelHash = channel.Hash
			user.TelegramChannelID = channel.TelegramChannelID
			user.TelegramChannelURL = channel.TelegramChannelURL
		}

		if err := s.deps.PG.CreateUser(user); err != nil {
			return s.InternalServerError(c, err)
		}
	}

	return s.ResponseOK(c)
}

type EventsSubmitLaunchRequest struct {
	BotToken   string `json:"bot_token"`
	TelegramID int64  `json:"telegram_id"`

	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}

func (req *EventsSubmitLaunchRequest) validate() error {
	if req.TelegramID == 0 {
		return errors.New("invalid telegram_id")
	}

	return nil
}

func (s *Server) EventsSubmitLaunchHandler(c *fiber.Ctx) error {
	req := &EventsSubmitLaunchRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	bot, err := s.deps.PG.SelectBotByToken(req.BotToken)
	if err != nil {
		return s.InternalServerError(c,
			fmt.Errorf("select bot by token: %w", err))
	}

	metricsEvents.WithLabelValues(bot.BotToken, types.EventTypeLaunch).Inc()

	users, err := s.deps.PG.SelectBotUsersByTelegramID(bot.ID, req.TelegramID)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	if len(users) > 0 {
		if err := s.deps.PG.UpdateUserMessagedAt(users[0].ID); err != nil {
			return s.InternalServerError(c, err)
		}

		ipInfo := net.ParseIP(req.IP)
		record, err := s.deps.GeoIP.Country(ipInfo)
		if err != nil {
			log.Warn("failed to procces ip", zap.Error(err))
		}

		uaInfo := uasurfer.Parse(req.UserAgent)

		if err := s.deps.PG.UpdateUserHeadersInfo(users[0].ID,
			&types.User{
				IP:          req.IP,
				UserAgent:   req.UserAgent,
				CountryCode: record.Country.IsoCode,
				OSName:      uaInfo.OS.Name.StringTrimPrefix(),
				DeviceType:  uaInfo.DeviceType.StringTrimPrefix(),
			}); err != nil {
			return s.InternalServerError(c, err)
		}
	}

	return s.ResponseOK(c)
}

type EventsSubmitDepositRequest struct {
	UserID             int   `json:"user_id"`
	ReporterTelegramID int64 `json:"reporter_telegram_id"`
}

type EventsSubmitDepositResponse struct {
	Success bool `json:"success"`
}

func (req *EventsSubmitDepositRequest) validate() error {
	if req.UserID == 0 {
		return errors.New("invalid user_id")
	}

	return nil
}

func (s *Server) EventsSubmitDepositHandler(c *fiber.Ctx) error {
	req := &EventsSubmitDepositRequest{}
	if err := c.BodyParser(&req); err != nil {
		return s.BadRequest(c, err)
	}

	if err := req.validate(); err != nil {
		return s.BadRequest(c, err)
	}

	el := &types.EventsLog{
		EventType:          types.EventTypeDeposit,
		ReporterTelegramID: req.ReporterTelegramID,
		UserID:             req.UserID,
	}

	if err := s.deps.PG.CreateEventLog(el); err != nil {
		return s.InternalServerError(c, err)
	}

	if err := s.deps.PG.UpdateUserDepositState(req.UserID); err != nil {
		return s.InternalServerError(c, err)
	}

	return s.ResponseOK(c)
}
