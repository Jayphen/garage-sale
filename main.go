package main

import (
	"fmt"
	"log"
	"net/http"

	"garagesale.jayphen.dev/assets/templ/layouts"
	"garagesale.jayphen.dev/assets/templ/pages"
	"garagesale.jayphen.dev/crontab"
	"garagesale.jayphen.dev/handlers"
	"garagesale.jayphen.dev/model"
	"garagesale.jayphen.dev/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	handlers.RegisterBidHandlers(app)
	handlers.RegisterSSEHandlers(app)
	crontab.RegisterCronJobs(app)

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", HomeHandler(e))
		e.Router.GET("/items", ItemsGet(e))

		e.Router.Static("/assets", "assets")

		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set("HX-Push-Url", "/")

		items, err := (&model.Item{}).GetItems(e.App.Dao())
		if err != nil {
			fmt.Println(err)
		}

		return utils.Render(c, http.StatusOK, layouts.Layout(items))
	}
}

func ItemsGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		items, err := (&model.Item{}).GetItems(e.App.Dao())
		if err != nil {
			fmt.Println(err)
		}

		return utils.Render(c, 200, pages.ItemsList(items))
	}
}
