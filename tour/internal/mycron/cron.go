package mycron

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func CronDemo()  {
	c := newWithSecond()

	c.AddFunc("@every 4s", func() {
		fmt.Println("每4秒执行一次")
	})

	c.Start()
	select {}
}

func newWithSecond() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
