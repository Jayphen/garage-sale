package crontab

import (
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
)

func RegisterCronJobs(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		scheduler := cron.New()

		resetPricing(app, scheduler)

		scheduler.Start()

		return nil
	})

	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		go func() {
			for {
				select {
				case <-done:
					ticker.Stop()
					return
				case <-ticker.C:
					decPricingTick(app)
				}
			}
		}()

		return nil
	})
}
