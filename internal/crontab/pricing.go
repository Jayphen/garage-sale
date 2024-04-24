package crontab

import (
	"fmt"
	"math"
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

func decPricingTick(app *pocketbase.PocketBase) error {
	const operationalEndHour = 21
	const operationalStartHour = 9
	const operationalHours = operationalEndHour - operationalStartHour

	currentTime := time.Now()

	// Exit if it is outside business hours
	if currentTime.Hour() < operationalStartHour || currentTime.Hour() >= operationalEndHour {
		fmt.Println("outside working hours")
		return nil
	}

	records, err := app.Dao().FindRecordsByFilter("items",
		"price > minPrice",
		"", 0, 0,
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Calculate when today's operational hours started
	todayStartOperational := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 9, 0, 0, 0, currentTime.Location())

	for _, r := range records {
		maxP := r.GetFloat("maxPrice")
		minP := r.GetFloat("minPrice")
		currentPrice := r.GetFloat("price")

		addedTimestamp := r.GetTime("created")

		// Adjust start time to today's 9 AM if the item was added before today
		if addedTimestamp.Before(todayStartOperational) {
			addedTimestamp = todayStartOperational
		}

		// Calculate the duration and left time in seconds for price decrease
		itemOperationalSeconds := float64((operationalEndHour - operationalStartHour) * 3600)
		elapsedSinceAdded := currentTime.Sub(addedTimestamp).Seconds()
		remainingSeconds := itemOperationalSeconds - elapsedSinceAdded

		if remainingSeconds <= 0 {
			r.Set("price", minP)
			app.Dao().SaveRecord(r)
			continue // Prevent division by zero and unnecessary adjustments after time has elapsed
		}

		totalDecrease := maxP - minP
		decreasePerSecond := totalDecrease / itemOperationalSeconds

		decreaseAmount := decreasePerSecond * 5 // every 5-second interval
		newPrice := currentPrice - decreaseAmount
		if newPrice < minP {
			newPrice = minP
		}

		// Round to two decimal places for monetary values
		newPriceRounded := math.Round(newPrice)

		r.Set("price", newPriceRounded)
		if err := app.Dao().SaveRecord(r); err != nil {
			return err // Properly handle potential errors in persistence
		}

		handlers.SendMessage("updated pricing")
	}

	return nil
}
