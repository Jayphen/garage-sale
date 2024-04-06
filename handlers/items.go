package handlers

import (
	"fmt"

	"garagesale.jayphen.dev/assets/templ/pages"
	"garagesale.jayphen.dev/model"
	"garagesale.jayphen.dev/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterItemsHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/items", ItemsGet(e))
		e.Router.GET("/items/:id/price", ItemPriceGet(e))
		return nil
	})
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

func ItemPriceGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")

		item, err := (&model.Item{}).FindItemById(e.App.Dao(), id)
		if err != nil {
			fmt.Println(err)
		}

		return utils.Render(c, 200, pages.Price(item))
	}
}
