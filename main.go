package main

import (
	"log"
	"net/http"

	"garagesale.jayphen.dev/assets/templ/layouts"
	"garagesale.jayphen.dev/assets/templ/pages"
	"garagesale.jayphen.dev/model"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

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

func Render(c echo.Context, statusCode int, t templ.Component) error {
	c.Response().Writer.WriteHeader(statusCode)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return t.Render(c.Request().Context(), c.Response().Writer)
}

func HomeHandler(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		c.Response().Header().Set("HX-Push-Url", "/")

		items, err := (&model.Item{}).GetItems(e.App.Dao())
		if err != nil {
			return err
		}

		return Render(c, http.StatusOK, layouts.Layout(items))
	}
}

func ItemsGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		items, err := (&model.Item{}).GetItems(e.App.Dao())
		if err != nil {
			return err
		}

		return Render(c, 200, pages.ItemsList(items))
	}
}
