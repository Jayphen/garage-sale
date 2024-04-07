package main

import (
	"log"

	"garagesale.jayphen.dev/crontab"
	"garagesale.jayphen.dev/handlers"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	handlers.RegisterHomeHandlers(app)
	handlers.RegisterBidHandlers(app)
	handlers.RegisterItemsHandlers(app)
	handlers.RegisterSSEHandlers(app)

	crontab.RegisterCronJobs(app)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.Static("/assets", "assets")

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
