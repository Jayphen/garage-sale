package handlers

import (
	"fmt"
	"time"

	"garagesale.jayphen.dev/internal/model"
	"garagesale.jayphen.dev/internal/utils"
	components "garagesale.jayphen.dev/ui/components/item"
	"garagesale.jayphen.dev/ui/pages"
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

		return utils.Render(c, 200, pages.ItemsList(items, true))
	}
}

func ItemGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")

		item, err := (&model.Item{Id: id}).FindItemById(e.App.Dao())
		if err != nil {
			return err
		}

		currentHour := time.Now().Hour()
		open := true
		if currentHour <= operationalStartHour || currentHour >= operationalEndHour {
			open = false
		}

		return utils.Render(c, 200, pages.ItemPage(item, utils.GetCartSize(c.Request()), open))
	}
}

func ItemPriceGet(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		id := c.PathParam("id")

		item, err := (&model.Item{Id: id}).FindItemById(e.App.Dao())
		if err != nil {
			return err
		}

		return utils.Render(c, 200, components.Price(components.ItemPrice{Id: id, Price: item.Price}))
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

		err := (&model.Item{Id: id}).SetItemStatus(e.App.Dao(), status)
		if err != nil {
			fmt.Println("Error setting item status:", err)
			return err
		}

		item, err := (&model.Item{Id: id}).FindItemById(e.App.Dao())
		if err != nil {
			return err
		}

		return utils.Render(c, 200, pages.ItemCard(item, true))
	}
}
