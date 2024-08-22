package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/cmd/worker/fbtool"
	"github.com/prosperofair/stata/pkg/pgsql"
	"github.com/prosperofair/stata/pkg/types"
	"go.uber.org/zap"
)

func (w *Worker) fbtoolFetcher() error {
	log.Info("fetching fbtool...")

	tokens, err := w.pg.SelectUnfetchedFBToolTokens()
	if err != nil {
		return fmt.Errorf("failed to select fbtool tokens: %w", err)
	}

	for _, token := range tokens {
		log.Info("fetching fbtool accounts", zap.String("token", token.Token))

		if err := w.fbtoolFetchAccounts(token); err != nil {
			return fmt.Errorf("failed to fetch fbtool accounts: %w", err)
		}

		if err := w.pg.UpdateFBToolTokenFetchedAt(token.ID); err != nil {
			return fmt.Errorf("failed to update fbtool token fetched_at: %w", err)
		}
		if err := w.pg.UpdateFBToolTokenDaysToFetch(token.ID); err != nil {
			return fmt.Errorf("failed to update days to fetch: %w", err)
		}
	}

	return nil
}

func (w *Worker) fbtoolFetchAccounts(token *types.FBToolToken) error {
	fc := fbtool.NewClient(token.Token)

	accounts, err := fc.GetAccounts()
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	for _, account := range accounts {
		if account.ID == 0 {
			continue
		}

		if err := w.pg.CreateFBToolAccount(
			&types.FBToolAccount{
				TokenID:           token.ID,
				FBToolAccountID:   account.ID,
				FBToolAccountName: account.Name,
			},
		); err != nil && !errors.Is(err, pgsql.ErrAlreadyExists) {
			log.Error("failed to create fbtool account", zap.Error(err))
			continue
		}
	}

	unfetched, err := w.pg.SelectUnfetchedFBToolAccounts(token.ID)
	if err != nil {
		return fmt.Errorf("failed to select unfetched fbtool accounts: %w", err)
	}

	for _, account := range unfetched {
		log.Info("fetching fbtool account stats", zap.Any("account", account))

		if err := w.fbtoolFetchAccountStats(fc, account, token.DaysToFetch); err != nil {
			log.Error("failed to fetch fbtool account stats", zap.Error(err))
			continue
		}
	}

	return nil
}

func (w *Worker) fbtoolFetchAccountStats(fc *fbtool.Client, account *types.FBToolAccount, daysToFetch int) error {
	now := time.Now()
	end := now.AddDate(0, 0, -1)
	start := now.AddDate(0, 0, -daysToFetch)

	fetched := false
	startedAt := time.Now()
	if err := w.pg.UpdateFBToolAccountFetchedAt(account.FBToolAccountID, fetched, 0); err != nil {
		return fmt.Errorf("failed to update fbtool account fetched_at: %w", err)
	}

	stats, err := fc.GetStatistics(account.FBToolAccountID, start, end)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %w", err)
	}

	for _, stat := range stats.Data {
		for _, ad := range stat.Ads.Data {

			if ad.Insights.Data != nil {
				for _, insight := range ad.Insights.Data {
					campaign := &types.FBToolCampaignStat{
						FBToolAccountID: account.FBToolAccountID,

						CampaignID:   ad.ID,
						CampaignName: ad.Name,

						Status:          ad.Status,
						EffectiveStatus: ad.EffectiveStatus,

						Impressions: insight.Impressions,
						Clicks:      insight.Clicks,
						Spend:       insight.Spend,
						Date:        insight.Date,
					}

					if err := w.pg.CreateFBToolCampaignStat(campaign); err != nil {
						log.Error("failed to create fbtool campaign stat", zap.Error(err))
						continue
					}
				}

				fetched = true
			}
		}
	}

	duration := int(time.Since(startedAt).Seconds())
	if err := w.pg.UpdateFBToolAccountFetchedAt(account.FBToolAccountID, fetched, duration); err != nil {
		return fmt.Errorf("failed to update fbtool account fetched_at: %w", err)
	}

	log.Info("fetched fbtool account stats", zap.Bool("fetched", fetched), zap.Int("duration", duration))

	return nil
}
