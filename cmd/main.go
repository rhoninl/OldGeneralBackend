package main

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	a, b := c.AddFunc("* * * * *", A)
	log.Println(a, b)
	c.Start()
	log.Println(time.Now())
	select {}
}

func A() {
	log.Println("++++++++++++++++++++++=")
	startZeroTime := time.UnixMicro(1781024975357991).Add(-8 * time.Hour)
	currentSigninNum := 1
	log.Println(startZeroTime.Add(24 * time.Hour * time.Duration(currentSigninNum)).Sub(time.Unix(1680969600, 0).Add(24 * time.Hour)))
}
