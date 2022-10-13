package usecron

import (
	"app/pkg/logging"
	"app/pkg/setting"
	"app/routers"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ehsaniara/gointerlock"
	"github.com/jasonlvhit/gocron"
)

func RunCron() {
	var (
		ctx    = context.Background()
		logger = logging.Logger{}
	)

	conString := fmt.Sprintf("%s:%d", setting.RedisDBSetting.Host, setting.RedisDBSetting.Port)
	jobTicker := gointerlock.GoInterval{
		Interval:      30 * time.Second,
		Arg:           myJob,
		Name:          "NotificationTicket",
		RedisHost:     conString,
		RedisPassword: setting.RedisDBSetting.Password,
	}

	err := jobTicker.Run(ctx)
	if err != nil {
		logger.Fatal("error cron job ", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		// cronJob.Stop()
		gocron.Clear()
		logger.Info("terminating: context cancelled")
	case <-sigterm:
		logger.Info("terminating: via signal")
	}
	wg.Wait()
}

func myJob() {
	var ctx = context.Background()
	routers.UseNotification.NotificationSms(ctx)
}

// func myJob() {
// 	fmt.Println(time.Now(), " - called")
// }
