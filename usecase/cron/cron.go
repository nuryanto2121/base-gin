package usecron

import (
	"app/routers"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jasonlvhit/gocron"
)

func RunCron() {
	var ctx = context.Background()
	// store.JobService.ProcessJob(ctx)
	gocron.Every(3).Minutes().From(gocron.NextTick()).Do(func() {
		// store.JobService.ProcessJob(ctx)
		// usenotification.NotificationSms(ctx)
		routers.UseNotification.NotificationSms(ctx)
	})
	<-gocron.Start()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		// cronJob.Stop()
		gocron.Clear()
		// log.Info("terminating: context cancelled")
	case <-sigterm:
		// log.Info("terminating: via signal")
	}
	wg.Wait()
}
