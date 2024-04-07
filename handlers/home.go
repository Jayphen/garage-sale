package handlers

import (
	"fmt"
	"net/http"

	"garagesale.jayphen.dev/assets/templ/layouts"
	"garagesale.jayphen.dev/model"
	"garagesale.jayphen.dev/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterHomeHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			c.Response().Header().Set("HX-Push-Url", "/")

			items, err := (&model.Item{}).GetItems(e.App.Dao())
			if err != nil {
				fmt.Println(err)
			}

			return utils.Render(c, http.StatusOK, layouts.Layout(items))
		})

		return nil
	})
}
