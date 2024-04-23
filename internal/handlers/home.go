package handlers

import (
	"fmt"
	"net/http"

	"garagesale.jayphen.dev/internal/model"
	"garagesale.jayphen.dev/internal/utils"
	"garagesale.jayphen.dev/ui/pages"
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

			cartSize := 0
			session := utils.GetSession(c.Request())

			if session != nil && session.Values["cart"] != nil {
				sessionCartSize, err := model.GetCartSize(e.App.Dao(), session.Values["cart"].(string))

				cartSize = sessionCartSize
				if err != nil {
					fmt.Println(err)
				}
			}

			return utils.Render(c, http.StatusOK, pages.ItemsPage(items, cartSize))
		})

		return nil
	})
}
