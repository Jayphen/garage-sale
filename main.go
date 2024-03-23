package main

import (
	"log"
	"net/http"

	"garagesale.jayphen.dev/assets/templ/layouts"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", HomeHandler)
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

func HomeHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Push-Url", "/")

	return Render(c, http.StatusOK, layouts.Layout())
}
