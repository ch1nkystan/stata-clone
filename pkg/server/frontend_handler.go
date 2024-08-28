package server

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prosperofair/stata/pkg/types"
)

type dateRangeRequest struct {
	BotToken string `json:"bot_token"`

	StartAt string    `json:"start_at"`
	Start   time.Time `json:"-"`

	EndAt string    `json:"end_at"`
	End   time.Time `json:"-"`

	StartPrev time.Time `json:"-"`
	EndPrev   time.Time `json:"-"`
}

func (req *dateRangeRequest) validate() error {
	if req.BotToken == "" {
		return errors.New("bot_token is empty")
	}

	sat, err := time.Parse(time.DateOnly, req.StartAt)
	if err != nil {
		sat = time.Now().Add(-time.Hour * 24 * 30)
	}

	req.Start = sat

	eat, err := time.Parse(time.DateOnly, req.EndAt)
	if err != nil {
		eat = time.Now()
	}

	if eat.After(time.Now().Add(time.Hour * 24)) {
		eat = time.Now().Add(time.Hour * 24)
	}

	req.End = eat

	// substract end from start to get the period
	diff := req.End.Sub(req.Start)

	if diff <= 0 {
		diff = time.Hour * 24
		req.End = req.Start.Add(diff)
	}

	req.StartPrev = req.Start.Add(-diff)
	req.EndPrev = req.End.Add(-diff)

	// convert to date only
	req.Start = req.Start.Truncate(time.Hour * 24)
	req.End = req.End.Truncate(time.Hour * 24)

	req.StartPrev = req.StartPrev.Truncate(time.Hour * 24)
	req.EndPrev = req.EndPrev.Truncate(time.Hour * 24)

	return nil
}

type conversionsByDayResponse struct {
	Data []*types.ConversionRow `json:"data"`
}

func (s *Server) conversionsByDayHandler(c *fiber.Ctx) error {
	req := &dateRangeRequest{}
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

	users, err := s.deps.PG.SelectBotUsersByDay(bot.ID, req.Start, req.End)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res := &conversionsByDayResponse{
		Data: make([]*types.ConversionRow, 0),
	}

	leads, err := s.deps.PG.SelectBotLeadsByDay(bot.ID, req.Start, req.End)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	// make same as above but from end to start
	for d := req.End; d.After(req.Start); d = d.Add(-time.Hour * 24) {
		row := &types.ConversionRow{
			ByDay:     d.Format(time.DateOnly),
			Deeplinks: make([]*types.ConversionRow, 0),
		}

		uData, ok := users[d.Format(time.DateOnly)]
		if ok {
			row.UsersTotal = uData.UsersTotal
			row.UsersUnique = uData.UsersUnique
			row.UsersUniqueRate = uData.UsersUniqueRate
			row.Deeplinks = uData.Deeplinks
		}

		lData, ok := leads[d.Format(time.DateOnly)]
		if ok {
			row.LeadsTotal = lData.LeadsTotal
			row.LeadsUsers = lData.LeadsUsers
			row.LeadsPerUser = lData.LeadsPerUser
			if row.UsersTotal != 0 {
				row.LeadsConversionRate = float64(row.LeadsUsers) / float64(row.UsersTotal) * 100
			}
			row.Profit = lData.Profit

			for _, d := range row.Deeplinks {
				dl, ok := lData.DeeplinksLeads[d.Label]
				if ok {
					d.LeadsTotal = dl.LeadsTotal
					d.LeadsUsers = dl.LeadsUsers
					d.LeadsPerUser = dl.LeadsPerUser
					if d.UsersTotal != 0 {
						d.LeadsConversionRate = float64(d.LeadsUsers) / float64(d.UsersTotal) * 100
					}

					d.Profit = dl.Profit
				}
			}
		}

		res.Data = append(res.Data, row)
	}

	return c.JSON(res)
}

func (s *Server) depositsLogHandler(c *fiber.Ctx) error {
	req := &dateRangeRequest{}
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

	deposits, err := s.deps.PG.SelectDepositsByBotID(bot.ID, req.Start, req.End)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	for i, d := range deposits {
		deposits[i].Hash = d.Hash[len(d.Hash)-10:]
	}

	res := &depositsLogResponse{
		Data: deposits,
	}

	return c.JSON(res)
}

type depositsLogResponse struct {
	Data []*types.DepositRow `json:"data"`
}

type metricsResponse struct {
	Data map[string]*types.MetricRow `json:"data"`

	Range *dateRange `json:"date_range"`
}

type dateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`

	StartPrev string `json:"start_prev"`
	EndPrev   string `json:"end_prev"`
}

func (s *Server) metricsHandler(c *fiber.Ctx) error {
	req := &dateRangeRequest{}
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

	res := &metricsResponse{
		Data: make(map[string]*types.MetricRow),

		Range: &dateRange{
			Start:     req.Start.Format(time.DateOnly),
			End:       req.End.Format(time.DateOnly),
			StartPrev: req.StartPrev.Format(time.DateOnly),
			EndPrev:   req.EndPrev.Format(time.DateOnly),
		},
	}

	users, err := s.deps.PG.SelectUsersMetric(bot.ID, req.Start, req.End, req.StartPrev, req.EndPrev)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res.Data["users"] = users

	usersUnique, err := s.deps.PG.SelectUsersUniqueMetric(bot.ID, req.Start, req.End, req.StartPrev, req.EndPrev)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res.Data["users_unique"] = usersUnique

	leads, err := s.deps.PG.SelectLeadsMetric(bot.ID, req.Start, req.End, req.StartPrev, req.EndPrev)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res.Data["leads"] = leads

	profit, err := s.deps.PG.SelectProfitMetric(bot.ID, req.Start, req.End, req.StartPrev, req.EndPrev)
	if err != nil {
		return s.InternalServerError(c, err)
	}

	res.Data["profit"] = profit

	// users, err := s.deps.PG.SelectBotUsersByDay(bot.ID, req.Start, req.End)

	// deposits, err := s.deps.PG.SelectDepositsByBotID(bot.ID, req.Start, req.End)
	// if err != nil {
	// 	return s.InternalServerError(c, err)
	// }

	// for i, d := range deposits {
	// 	deposits[i].Hash = d.Hash[len(d.Hash)-10:]
	// }

	// res := &depositsLogResponse{
	// 	Data: deposits,
	// }

	return c.JSON(res)
}
