package main

import (
	"fmt"
	"time"

	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/cmd/worker/fbtool"
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

		if err := w.pg.UpdateFBToolTokenFetchedAt(token.ID); err != nil {
			return fmt.Errorf("failed to update fbtool token fetched_at: %w", err)
		}

		if err := w.fbtoolFetchAccounts(token); err != nil {
			return fmt.Errorf("failed to fetch fbtool accounts: %w", err)
		}
	}

	return nil
}

func (w *Worker) fbtoolFetchAccounts(token *types.FBToolToken) error {
	fc := fbtool.NewClient(token.Token)

	time.Sleep(5 * time.Second)
	accounts, err := fc.GetAccounts()
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	for _, account := range accounts {
		log.Info("fetching fbtool account", zap.Any("account", account))
		if account.ID == 0 {
			log.Error("invalid account id, skipping", zap.Any("account", account))
			continue
		}

		current, err := w.pg.SelectFBToolAccountsByAccountID(account.ID)
		if err != nil {
			log.Error("failed to select fbtool account", zap.Error(err))
			continue
		}

		if len(current) == 0 {
			if err := w.pg.CreateFBToolAccount(&types.FBToolAccount{
				TokenID:           token.ID,
				FBToolAccountID:   account.ID,
				FBToolAccountName: account.Name,
			}); err != nil {
				log.Error("failed to create fbtool account", zap.Error(err))
				continue
			}
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
	start := time.Now().AddDate(0, 0, -1)

	fetched := false
	nodata := 0

	for i := 0; i < daysToFetch; i++ {
		time.Sleep(3 * time.Second)

		if nodata > 3 {
			log.Debug("no data for 3 days, disabling and skipping", zap.Int("account_id", account.FBToolAccountID))
			if err := w.pg.DisableFBToolAccount(account.FBToolAccountID); err != nil {
				return fmt.Errorf("failed to disable fbtool account: %w", err)
			}
			return nil
		}

		date := start.AddDate(0, 0, -i)

		log.Debug("fetching fbtool account stats",
			zap.String("api_key", fc.APIKey),
			zap.Int("account_id", account.FBToolAccountID),
			zap.String("account_name", account.FBToolAccountName),
			zap.String("date", date.Format("2006-01-02")),
		)

		stats, err := fc.GetStatistics(account.FBToolAccountID, date)
		if err != nil {
			log.Error("failed to get statistics", zap.Error(err))
			continue
		}

		if stats == nil || stats.Data == nil {
			nodata++
			log.Debug("no data for date",
				zap.String("api_key", fc.APIKey),
				zap.Int("account_id", account.FBToolAccountID),
				zap.String("account_name", account.FBToolAccountName),
				zap.String("date", date.Format("2006-01-02")),
				zap.Any("response", stats),
			)
			continue
		}

		for _, stat := range stats.Data {
			for _, ad := range stat.Ads.Data {
				campaign := &types.FbToolCampaignStat{
					FBToolAccountID: account.FBToolAccountID,
					CampaignID:      ad.ID,
					CampaignName:    ad.Name,
					Status:          ad.Status,
					EffectiveStatus: ad.EffectiveStatus,
				}

				if ad.Insights.Data != nil {
					for _, insight := range ad.Insights.Data {
						campaign.Impressions += insight.Impressions
						campaign.Clicks += insight.Clicks
						campaign.Spend += insight.Spend
					}
				}

				campaign.Date = tsod(date)

				if err := w.pg.CreateFBToolCampaignStat(campaign); err != nil {
					log.Error("failed to create fbtool campaign stat", zap.Error(err))
					continue
				}

				fetched = true
				nodata = 0
			}
		}
	}

	if fetched {
		if err := w.pg.UpdateFBToolAccountFetchedAt(account.FBToolAccountID); err != nil {
			return fmt.Errorf("failed to update fbtool account fetched_at: %w", err)
		}
	}

	return nil
}

func tsod(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}
