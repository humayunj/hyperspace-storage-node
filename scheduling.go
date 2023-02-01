package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func InitSchedular() {
	s := gocron.NewScheduler(time.UTC)
	t := 0
	tt := time.Now().Add(5 * time.Second).UTC()
	printLn(tt)
	s.Every(1).StartAt(tt).LimitRunsTo(1).Do(func() {
		printLn("Sch: " + fmt.Sprintf("%v", time.Now().UTC().String()))
		t++
	})
	s.StartAsync()
}
