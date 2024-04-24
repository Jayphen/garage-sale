package main

import (
	"log"
	"time"

	"garagesale.jayphen.dev/internal/crontab"
	"garagesale.jayphen.dev/internal/handlers"
	"garagesale.jayphen.dev/internal/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func init() {
	time.LoadLocation("Europe/Vienna")

	utils.CreateStore()
}

func main() {
	app := pocketbase.New()

	handlers.RegisterHomeHandlers(app)
	handlers.RegisterItemsHandlers(app)
	handlers.RegisterCheckoutHandlers(app)
	handlers.RegisterCartHandlers(app)
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
