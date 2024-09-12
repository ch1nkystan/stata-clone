package server

import (
	"time"

	"github.com/prosperofair/stata/pkg/types"
)

const (
	GroupByDay   = "day"
	GroupByWeek  = "week"
	GroupByMonth = "month"
)

// calcPeriod - вычисляет начало и конец периода
func calcPeriod(groupBy string, start, end time.Time) (time.Time, time.Time) {
	var ps, pe time.Time

	switch groupBy {
	case GroupByDay:
		ps = start
		pe = end
	
	// Находим начало конец первой и начало последней недели, а затем итерируемся по неделям
	// Пример: на входе 2024-06-28, 2024-09-05; на выходе 2024-06-23, 2024-09-02
	case GroupByWeek:
		ps = start.AddDate(0, 0, -int(start.Weekday()))
		pe = end.AddDate(0, 0, -int(end.Weekday())+1)
	
	// Находим первый день первого и последнего месяца, а затем итерируемся по месяцам
	// Пример: на входе 2024-06-30, 2024-09-05; на выходе 2024-06-01, 2024-09-01
	case GroupByMonth:
		ps = time.Date(start.Year(), start.Month(), 0, 0, 0, 0, 0, start.Location())
		pe = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())
	}

	return ps, pe
}

// сalcPeriodEnd - вычисляет конец периода для конкретной итерации
func сalcPeriodEnd(d time.Time, period string) string {
	switch period {
	// Находим последний день недели
	// Пример: на входе 2024-06-24; на выходе 2024-06-30
	case GroupByWeek:
		return d.AddDate(0, 0, 6).Format(time.DateOnly)
	
	// Находим последний день месяца
	// Пример: на входе 2024-06-01; на выходе 2024-06-30
	case GroupByMonth:
		return d.AddDate(0, 1, -1).Format(time.DateOnly)
	}

	return d.Format(time.DateOnly)
}

func decrement(period string) func(time.Time) time.Time {
	switch period {
	case GroupByDay:
		return func(t time.Time) time.Time { return t.AddDate(0, 0, -1) } // Итерируемся по дням
	case GroupByWeek:
		return func(t time.Time) time.Time { return t.AddDate(0, 0, -7) } // Итерируемся по неделям
	case GroupByMonth:
		return func(t time.Time) time.Time { return t.AddDate(0, -1, 0) } // Итерируемся по месяцам
	default:
		return nil
	}
}

type periodMetricsConfig struct {
	users    map[string]*types.ConversionRow
	leads    map[string]*types.ConversionRow
	expenses map[string]*types.ConversionRow

	groupBy string
	start   time.Time
	end     time.Time
}

func calcPeriodMetrics(cfg *periodMetricsConfig) []*types.ConversionRow {
	res := make([]*types.ConversionRow, 0)

	dec := decrement(cfg.groupBy)
	if dec == nil {
		return res
	}

	start, end := calcPeriod(cfg.groupBy, cfg.start, cfg.end)

	for d := end; d.After(start); d = dec(d) {
		row := &types.ConversionRow{
			PeriodStart: d.Format(time.DateOnly),
			PeriodEnd:   сalcPeriodEnd(d, cfg.groupBy),
		}

		if uData, ok := cfg.users[d.Format(time.DateOnly)]; ok {
			row.UsersTotal = uData.UsersTotal
			row.UsersUnique = uData.UsersUnique
			row.UsersUniqueRate = uData.UsersUniqueRate
		}

		if lData, ok := cfg.leads[d.Format(time.DateOnly)]; ok {
			row.LeadsTotal = lData.LeadsTotal
			row.LeadsUsers = lData.LeadsUsers
			row.LeadsPerUser = lData.LeadsPerUser
			if row.UsersTotal != 0 {
				row.LeadsConversionRate = float64(row.LeadsUsers) / float64(row.UsersTotal) * 100
			}

			row.Income = lData.Income
		}

		if eData, ok := cfg.expenses[d.Format(time.DateOnly)]; ok {
			row.Impressions = eData.Impressions
			row.Clicks = eData.Clicks
			row.Expense = eData.Expense
		}

		row.Profit = row.Income - row.Expense

		res = append(res, row)
	}

	return res
}
