package main

import (
	"github.com/brianvoe/gofakeit/v7"
	"time"
)

func RandomTimeBetween(start, end time.Time) time.Time {
	if start.After(end) {
		start, end = end, start
	}
	delta := end.Sub(start)
	randomDuration := time.Duration(gofakeit.Float64Range(0, float64(delta)))
	return start.Add(randomDuration)
}
