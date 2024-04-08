package handlers

import (
	"fmt"

	"garagesale.jayphen.dev/model"
	components "garagesale.jayphen.dev/ui/components/item"
	"garagesale.jayphen.dev/ui/pages"
	"garagesale.jayphen.dev/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterItemsHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/items", ItemsGet(e, app))
		e.Router.GET("/items/:id", ItemGet(e))
		e.Router.GET("/items/:id/price", ItemPriceGet(e))
		e.Router.POST("/items/:id/status", ItemStatusSet(e))
		return nil
	})
}

func ItemsGet(e *core.ServeEvent, app *pocketbase.PocketBase) func(echo.Context) error {
	return func(c echo.Context) error {
		items, err := (&model.Item{}).GetItems(e.App.Dao())
		if err != nil {
			return err
		}

		return utils.Render(c, 200, pages.ItemsList(items))
	}
}

func ItemGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")

		item, err := (&model.Item{}).FindItemById(e.App.Dao(), id)
		if err != nil {
			return err
		}

		return utils.Render(c, 200, pages.ItemPage(item))
	}
}

func ItemPriceGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")

		item, err := (&model.Item{}).FindItemById(e.App.Dao(), id)
		if err != nil {
			return err
		}

		return utils.Render(c, 200, components.Price(item))
	}
}

func ItemStatusSet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		var status model.ItemStatus

		// get id from param and new status from POST body
		id := c.PathParam("id")

		if err := status.ParseFormValue(c.FormValue("status")); err != nil {
			fmt.Println("Error parsing form value:", err)
			return err
		}

		err := (&model.Item{}).SetItemStatus(e.App.Dao(), id, status)
		if err != nil {
			fmt.Println("Error setting item status:", err)
			return err
		}

		item, err := (&model.Item{}).FindItemById(e.App.Dao(), id)
		if err != nil {
			return err
		}

		return utils.Render(c, 200, pages.ItemCard(item))
	}
}
