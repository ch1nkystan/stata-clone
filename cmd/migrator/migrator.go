package main

import (
	"time"

	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"go.uber.org/zap"
)

type Migrator struct {
	old       *Client
	new       *pgsql.Client
	chunkSize int
	dryRun    bool
}

func NewMigrator(new *pgsql.Client, old *Client, cz int, dryRun bool) *Migrator {
	return &Migrator{
		new:       new,
		old:       old,
		chunkSize: cz,
		dryRun:    dryRun,
	}
}

func (m *Migrator) Run() error {
	users, err := m.new.SelectAllUsersWithEmptyCreationEvent()
	if err != nil {
		return err
	}

	notFound := 0

	bots, err := m.new.SelectAllBots()
	if err != nil {
		return err
	}

	for _, user := range users {
		bot, ok := bots[user.BotID]
		if !ok {
			log.Error("bot not found", zap.Int64("telegram_id", user.TelegramID), zap.Int("bot_id", user.BotID))
			continue
		}

		incomes, err := m.old.SelectIncomesInfoForUser(user.TelegramID, bot.BotUsername)
		if err != nil {
			return err
		}

		log.Info("incomes found", zap.Int("count", len(incomes)), zap.Int64("telegram_id", user.TelegramID))
		if len(incomes) == 0 {
			notFound++
			continue
		}

		income := incomes[0]
		deeplinks, err := m.new.SelectBotDeeplinksByLabel(bot.ID, income.IncomeSource)
		if err != nil {
			return err
		}

		if len(deeplinks) == 0 {
			log.Error("deeplink not found", zap.String("label", income.IncomeSource), zap.Int64("telegram_id", user.TelegramID))

			if income.IncomeSource == "bot" {
				if err := m.new.UpdateUserDeeplinkAndCreatedAt(user.ID, 10, income.CreateAt); err != nil {
					return err
				}
				log.Info("deeplink updated", zap.Int64("telegram_id", user.TelegramID), zap.Int("deeplink_id", 10), zap.String("label", income.IncomeSource), zap.Time("created_at", income.CreateAt))
			}

			continue
		}

		deeplink := deeplinks[0]

		if err := m.new.UpdateUserDeeplinkAndCreatedAt(user.ID, deeplink.ID, income.CreateAt); err != nil {
			return err
		}

		log.Info("deeplink updated", zap.Int64("telegram_id", user.TelegramID), zap.Int("deeplink_id", deeplink.ID), zap.String("label", income.IncomeSource), zap.Time("created_at", income.CreateAt))
	}

	time.Sleep(10 * time.Minute)

	return nil
}
