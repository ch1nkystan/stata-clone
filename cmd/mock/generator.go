package main

import (
	"errors"
	"fmt"
	"github.com/avct/uasurfer"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/prosperofair/pkg/log"
	"github.com/prosperofair/stata/pkg/pgsql"
	"github.com/prosperofair/stata/pkg/types"
	"go.uber.org/zap"
	"time"
)

var devices = []uasurfer.DeviceType{
	uasurfer.DeviceUnknown,
	uasurfer.DeviceComputer,
	uasurfer.DeviceTablet,
	uasurfer.DevicePhone,
	uasurfer.DeviceConsole,
	uasurfer.DeviceWearable,
	uasurfer.DeviceTV,
}

var os = []uasurfer.OSName{
	uasurfer.OSUnknown,
	uasurfer.OSWindowsPhone,
	uasurfer.OSWindows,
	uasurfer.OSMacOSX,
	uasurfer.OSiOS,
	uasurfer.OSAndroid,
	uasurfer.OSBlackberry,
	uasurfer.OSChromeOS,
	uasurfer.OSKindle,
	uasurfer.OSWebOS,
	uasurfer.OSLinux,
	uasurfer.OSPlaystation,
	uasurfer.OSXbox,
	uasurfer.OSNintendo,
	uasurfer.OSBot,
}

// randomSample create a data sample containing a percentage of a slice element in random order
func randomSample[T any](slice []T, percentage float64) []T {
	if len(slice) == 0 {
		return nil
	}

	count := int(float64(len(slice)) * percentage)
	if count == 0 {
		count = 1
	}
	if count > len(slice) {
		count = len(slice)
	}

	indices := make([]int, 0, len(slice))
	for i := range slice {
		indices = append(indices, i)
	}
	gofakeit.ShuffleAnySlice(&indices)

	sample := make([]T, 0, count)
	for _, v := range indices[:count] {
		sample = append(sample, slice[v])
	}

	return sample
}

type Generator struct {
	client                 *pgsql.Client
	minUserPerDay          int
	maxUserPerDay          int
	deepLinksCount         int
	referralUserPercentage float64
	referralUserCountMin   int
	referralUserCountMax   int
	period                 time.Duration
	botID                  int

	languages     []string
	events        []string
	mailingStates []string
	devicesString []string
	osString      []string
	labels        []string

	generatedTelegramId map[int64]bool

	logger *log.Logger
}

// generateUniqueTelegramID generates a unique Telegram user ID
// The ID is randomly selected between 100,000,000 and 999,999,999,999 — a range real-world Telegram user probably won't have IDs in this range
func (g *Generator) generateUniqueTelegramID() int64 {
	for {
		telegramID := int64(gofakeit.Float64Range(100_000_000, 999_999_999_999))
		if !g.generatedTelegramId[telegramID] {
			g.generatedTelegramId[telegramID] = true
			return telegramID
		}
	}
}

func newGenerator(
	client *pgsql.Client,
	minUsersPerDay,
	maxUsersPerDay,
	deepLinksCount int,
	referralUserPercentage float64,
	referralUserCountMin int,
	referralUserCountMax int,
	periodDays time.Duration,
	botID int,
	logger *log.Logger,
) (*Generator, error) {

	if minUsersPerDay < 1 {
		return nil, fmt.Errorf("minUsersPerDay must be >= 1")
	}
	if maxUsersPerDay < minUsersPerDay {
		return nil, fmt.Errorf("maxUsersPerDay must be more >= minUsersPerDay")
	}
	if deepLinksCount < 1 {
		return nil, fmt.Errorf("deepLinkCount must be >= 1")
	}
	if referralUserPercentage < 0 || referralUserPercentage > 1 {
		return nil, fmt.Errorf("referralUserPercentage must be between 0 and 1")
	}
	if referralUserCountMin < 1 {
		return nil, fmt.Errorf("referralUserCountMin must be >= 1")
	}
	if referralUserCountMax < referralUserCountMin {
		return nil, fmt.Errorf("referralUserCountMax must be >= referralUserCountMin")
	}
	if periodDays < 0 {
		return nil, fmt.Errorf("periodDays must be >= 0")
	}

	generator := &Generator{
		client:                 client,
		minUserPerDay:          minUsersPerDay,
		maxUserPerDay:          maxUsersPerDay,
		deepLinksCount:         deepLinksCount,
		referralUserPercentage: referralUserPercentage,
		referralUserCountMin:   referralUserCountMin,
		referralUserCountMax:   referralUserCountMax,
		period:                 periodDays,
		botID:                  botID,
		logger:                 logger,
	}

	generator.languages = []string{"ru", "en", "da", "de", "es"} // maybe need to add more
	generator.events = []string{types.EventTypeRegister, types.EventTypeMessage}
	generator.mailingStates = types.MailingStages

	strings := make([]string, 0, len(devices))
	for _, device := range devices {
		strings = append(strings, device.String())
	}
	generator.devicesString = strings

	strings = make([]string, 0, len(os))
	for _, os := range os {
		strings = append(strings, os.String())
	}
	generator.osString = strings

	generator.generatedTelegramId = make(map[int64]bool)

	generator.labels = []string{types.DeeplinkLabelReferral}

	return generator, nil
}

// fakeUser generate fake user
func (g *Generator) fakeUser(date time.Time, bot *types.Bot, deeplink *types.Deeplink) types.User {
	date = RandomTimeBetween(date.Add(-time.Hour*6), date.Add(time.Hour*6))

	user := types.User{
		TelegramID:        g.generateUniqueTelegramID(), // rarely, but can be collisions
		Firstname:         gofakeit.FirstName(),
		Lastname:          gofakeit.LastName(),
		ForwardSenderName: gofakeit.Name(),
		Username:          gofakeit.Username(),

		IP:          gofakeit.IPv4Address(),
		UserAgent:   gofakeit.UserAgent(),
		CountryCode: gofakeit.CountryAbr(),
		OSName:      gofakeit.RandomString(g.osString),
		DeviceType:  gofakeit.RandomString(g.devicesString),

		IsBot:     gofakeit.Bool(),
		IsPremium: gofakeit.Bool(),

		LanguageCode: gofakeit.RandomString(g.languages),

		DepositsTotal: gofakeit.Number(0, 10),
		DepositsSum:   gofakeit.Float64Range(10.0, 10000.0),
		Deposited:     gofakeit.Bool(),
		DepositedAt:   date,

		Seen: gofakeit.Number(0, 1000),

		EventCreated: gofakeit.RandomString(g.events),

		MailingState:          gofakeit.RandomString(g.mailingStates),
		MailingFailedAttempts: gofakeit.Number(0, types.UserMailingMaxFailedAttempts),
		MailingStateUpdatedAt: date,

		MessagedAt: date,
		CreatedAt:  date,
		UpdatedAt:  date,

		Subscribed:     gofakeit.Bool(),
		SubscribedAt:   date,
		UnsubscribedAt: date,
	}

	if deeplink != nil {
		user.DeeplinkID = deeplink.ID
	}

	if bot != nil {
		user.BotID = bot.ID
	}

	return user
}

// fakeDeepLink generate fake deepLink
func (g *Generator) fakeDeepLink(date time.Time, bot *types.Bot, referralUser *types.User) types.Deeplink {
	deepLink := types.Deeplink{
		Hash:      types.GenerateDeeplinkHash(),
		Active:    gofakeit.Bool(),
		Label:     gofakeit.RandomString(g.labels),
		CreatedAt: RandomTimeBetween(date, date.Add(time.Hour)),
		UpdatedAt: RandomTimeBetween(date, date.Add(time.Hour)),
	}

	if bot != nil {
		deepLink.BotID = bot.ID
	}

	if referralUser != nil {
		deepLink.ReferralTelegramID = referralUser.TelegramID
	}

	return deepLink
}

func (g *Generator) Generate() error {
	var (
		now  = time.Now()
		past = now.Add(-g.period)
	)

	bot, err := g.client.SelectBotByID(g.botID)
	if err != nil {
		return fmt.Errorf("error selecting bot by ID: %s", err)
	}
	if bot == nil {
		return errors.New("bot not found")
	}
	g.logger.Info("Bot selected", zap.Int("ID", bot.ID))

	g.logger.Info("Creating ad deeplinks", zap.Int("count", g.deepLinksCount))

	var adDeepLinks []types.Deeplink
	for range g.deepLinksCount {
		adDeepLinks = append(adDeepLinks, g.fakeDeepLink(past, bot, nil))
	}

	// rarely, but hash collisions are possible so we ignore them
	adDeepLinks, err = g.client.CreateDeeplinksIgnoreInsertErrorsReturningAll(adDeepLinks)
	if err != nil {
		return fmt.Errorf("error creating adDeepLinks: %s", err)
	}

	g.logger.Info("ad deeplinks created", zap.Int("count", len(adDeepLinks)))

	g.logger.Info("start creating ad users")
	var adUserTelegramIDs []int64
	var adUsers []types.User
	for date := past; !date.After(now); date = date.AddDate(0, 0, 1) {
		userCount := gofakeit.IntRange(g.minUserPerDay, g.maxUserPerDay)

		g.logger.Info("creating ad user", zap.Int("count", userCount), zap.Time("date", date))

		counts := make([]int, len(adDeepLinks))
		for range userCount {
			idx := gofakeit.IntRange(0, len(adDeepLinks)-1)
			counts[idx]++
		}

		for i, count := range counts {
			if count == 0 {
				continue
			}

			for range count {
				user := g.fakeUser(date, bot, &adDeepLinks[i])
				adUsers = append(adUsers, user)
				adUserTelegramIDs = append(adUserTelegramIDs, user.TelegramID)
			}
		}
	}

	if err = g.client.CreateUsers(adUsers); err != nil {
		return fmt.Errorf("failed to create ad users for deeplink: %w", err)
	}
	g.logger.Info("ad users created", zap.Int("count", len(adUsers)))

	// for percentage of users create referral adDeepLinks
	sample := randomSample(adUserTelegramIDs, g.referralUserPercentage)
	g.logger.Info("start creating referral deeplinks", zap.Int("count", len(sample)))
	var referralDeepLinks []types.Deeplink
	for _, telegramID := range sample {
		referralDeepLinks = append(referralDeepLinks,
			g.fakeDeepLink(
				RandomTimeBetween(past, now),
				bot,
				&types.User{TelegramID: telegramID},
			),
		)
	}

	referralDeepLinks, err = g.client.CreateDeeplinksIgnoreInsertErrorsReturningAll(referralDeepLinks)
	if err != nil {
		return fmt.Errorf("failed to create referral deeplinks: %w", err)
	}
	g.logger.Info("referral deeplinks created", zap.Int("count", len(referralDeepLinks)))

	// for each referral deeplink creating referralUserCountMin < count < referralUserCountMax referral users
	g.logger.Info("start creating referral users for referral deeplinks")
	var referralUsers []types.User
	for _, referralDeepLink := range referralDeepLinks {
		referralUserCount := gofakeit.IntRange(g.referralUserCountMin, g.referralUserCountMax)

		for range referralUserCount {
			referralUsers = append(referralUsers, g.fakeUser(RandomTimeBetween(referralDeepLink.CreatedAt, now), bot, &referralDeepLink))
		}
	}
	g.logger.Info("referral users created", zap.Int("count", len(referralUsers)))

	if err = g.client.CreateUsers(referralUsers); err != nil {
		return fmt.Errorf("failed to create referral users: %w", err)
	}

	g.logger.Info("stats",
		zap.Int("totalUsers", len(adUsers)+len(referralUsers)),
		zap.Int("totalDeepLinks", len(adDeepLinks)+len(referralDeepLinks)),
	)

	return nil
}
