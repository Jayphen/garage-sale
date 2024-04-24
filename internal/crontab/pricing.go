package crontab

import (
	"fmt"
	"time"

	"garagesale.jayphen.dev/internal/handlers"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/cron"
)

func resetPricing(app *pocketbase.PocketBase, scheduler *cron.Cron) {
	frequency := "10 21 * * *"

	scheduler.MustAdd("pricingReset", frequency, func() {
		app.Dao().DB().NewQuery("UPDATE items SET price = maxPrice").
			Execute()
	})
}

func decPricingTick() error {
	const operationalEndHour = 21
	const operationalStartHour = 9
	const operationalHours = operationalEndHour - operationalStartHour

	currentTime := time.Now()

	// Exit if it is outside business hours
	if currentTime.Hour() < operationalStartHour || currentTime.Hour() >= operationalEndHour {
		fmt.Println("outside working hours")
		return nil
	}

	handlers.SendMessage("updated pricing")

	return nil
}
