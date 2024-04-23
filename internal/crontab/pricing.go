package crontab

import (
	"fmt"
	"math"

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

func decPricingTick(app *pocketbase.PocketBase) error {
	records, err := app.Dao().FindRecordsByFilter("items",
		"price > minPrice",
		"", 0, 0,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, r := range records {
		p := r.GetFloat("price")
		maxP := r.GetFloat("maxPrice")
		minP := r.GetFloat("minPrice")
		roundedPrice := int(math.Round(p - ((maxP-minP)/8640)*100/100))

		r.Set("price", roundedPrice)

		app.Dao().SaveRecord(r)
	}
	handlers.SendMessage("updated pricing")

	return nil
}
