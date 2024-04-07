package crontab

import (
	"fmt"

	"garagesale.jayphen.dev/handlers"
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

func decPricing(app *pocketbase.PocketBase, scheduler *cron.Cron) {
	frequency := "* * * * *" // every minute

	scheduler.MustAdd("pricingDec", frequency, func() {
		_, err := app.Dao().DB().NewQuery("UPDATE items SET price = price - ((maxPrice - minPrice) / 720) WHERE price > minPrice").
			Execute()
		if err != nil {
			fmt.Println(err)
		} else {
			handlers.SendMessage("updated pricing")
		}
	})
}
