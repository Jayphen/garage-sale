package crontab

import (
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/cron"
)

const (
	operationalStartHour = 7
	operationalEndHour   = 14
)

func RegisterCronJobs(app *pocketbase.PocketBase) {
	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
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
					currentHour := time.Now().Hour()
					if currentHour >= operationalStartHour && currentHour < operationalEndHour {
						err := decPricingTick()
						if err != nil {
							fmt.Println("Error executing pricing tick:", err)
						}
					}
				}
			}
		}()

		return nil
	})
}
