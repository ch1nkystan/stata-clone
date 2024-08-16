package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/prosperofair/pkg/depot"
	"github.com/prosperofair/pkg/log"

	"github.com/prosperofair/stata/pkg/pgsql"
)

type Server struct {
	App *fiber.App

	cfg  *Config
	deps *Deps
}

type Config struct {
	ExposeMetrics bool

	BackendTokens  map[string]struct{}
	FrontendTokens map[string]struct{}
}

type Deps struct {
	PG *pgsql.PGSQLClient

	Depot *depot.Client
}

func New(cfg *Config, deps *Deps) *Server {
	s := &Server{
		App: fiber.New(fiber.Config{}),

		cfg:  cfg,
		deps: deps,
	}

	s.App.Use(cors.New())

	// monitoring
	s.App.Use(s.MonitoringMiddleware)

	// panic recovery
	s.App.Use(recover.New(recover.Config{
		Next:             nil,
		EnableStackTrace: true,
	}))

	if cfg.ExposeMetrics {
		log.Info("run metrics exposer...")
		go runMetricsExposer()
	}

	// status checks
	s.App.Get("/healthz", s.healthzHandler)

	api := s.App.Group("/api", s.apiMiddlewareBackend)

	bot := api.Group("/bots")
	bot.Post("/register", s.BotsRegisterHandler)
	bot.Post("/import", s.BotsImportHandler)

	user := api.Group("/users")
	user.Post("/search", s.UsersSearchHandler)
	user.Post("/get", s.UsersGetHandler)
	user.Post("/set/default-telegram-channel", s.UsersSetDefaultChannelHandler)
	user.Post("/update/telegram-channel", s.UsersUpdateTelegramChannelHandler)

	deeplink := api.Group("/deeplinks")
	deeplink.Post("/create", s.DeeplinksCreateHandler)
	deeplink.Post("/list", s.DeeplinksListHandler)
	deeplink.Post("/update", s.DeeplinksUpdateHandler)

	event := api.Group("/events")
	event.Post("/submit/user-register", s.EventsSubmitUserRegisterHandler)
	event.Post("/submit/message", s.EventsSubmitMessageHandler)
	event.Post("/submit/deposit", s.EventsSubmitDepositHandler)

	mailing := api.Group("/mailing")
	mailing.Post("/prepare/users-list", s.MailingPrepareUsersListHandler)
	mailing.Post("/finish/users-list", s.MailingFinishUsersListHandler)
	mailing.Post("/update/user-state", s.MailingUpdateUserStateHandler)

	transactions := api.Group("/transactions")
	transactions.Post("/create", s.TransactionsCreateHandler)

	stats := api.Group("/stats")
	stats.Post("/mailing-state", s.StatsMailingStateHandler)
	stats.Post("/users-count", s.StatsUsersCountHandler) // used by depot to get bots stats

	// method used by frontend to get stats
	f := s.App.Group("/f/api", s.apiMiddlewareFrontend)
	f.Post("/stats/conversions-by-day", s.conversionsByDayHandler)
	f.Post("/stats/deposits-log", s.depositsLogHandler)

	return s
}

func (s *Server) healthzHandler(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func runMetricsExposer() {
	const (
		timeout = 5 * time.Second
		addr    = ":2112"
	)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		srv := &http.Server{
			Addr:    addr,
			Handler: http.TimeoutHandler(mux, timeout, "request timeout"),
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("failed to run metrics exposer")
		}
	}()
}

type response struct {
	Message string `json:"message"`
}

func (s *Server) ResponseOK(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(response{"ok"})
}

func (s *Server) InternalServerError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusInternalServerError).
		JSON(response{err.Error()})
}

func (s *Server) BadRequest(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusBadRequest).
		JSON(response{err.Error()})
}

func (s *Server) Unauthorized(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).
		JSON(response{fiber.ErrUnauthorized.Message})
}

func (s *Server) MonitoringMiddleware(c *fiber.Ctx) error {
	start := time.Now()
	err := c.Next()
	finish := time.Now()
	latency := finish.Sub(start)

	log.Info("request handled",
		zap.Int("status", c.Response().StatusCode()),
		zap.Duration("latency", latency),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("url", c.OriginalURL()),
		zap.String("headers_raw", c.Request().Header.String()),
		zap.Error(err),
	)

	// status codes metrics
	metricStatusCodes.WithLabelValues(strconv.Itoa(c.Response().StatusCode())).Inc()

	return err
}
