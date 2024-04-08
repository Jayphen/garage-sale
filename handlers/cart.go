package handlers

import (
	"fmt"

	"garagesale.jayphen.dev/model"
	toast "garagesale.jayphen.dev/ui/components"
	components "garagesale.jayphen.dev/ui/components/cart"
	"garagesale.jayphen.dev/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func addToCart(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		itemId := c.FormValue("itemId")

		item, err := (&model.Item{}).FindItemById(e.App.Dao(), itemId)
		if err != nil {
			fmt.Println(err)
			return err
		}

		session := utils.GetSession(c.Request())

		// Retrieve cart ID from the session if it exists
		cartId, ok := session.Values["cart"].(string)
		if !ok {
			cartId = ""
		}

		newCartId, err := model.AddToCart(e.App.Dao(), item.Id, cartId)
		if err != nil {
			fmt.Println(err)
			return utils.Render(c, 400, toast.Toast("That's already in your cart. You should probably hurry up."))
		}

		session.Values["cart"] = newCartId
		session.Save(c.Request(), c.Response())

		cartSize, err := model.GetCartSize(e.App.Dao(), newCartId)
		if err != nil {
			fmt.Println(err)
		}

		return utils.Render(c, 200, components.Indicator(cartSize))
	}
}

func RegisterCartHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/cart", addToCart(e))
		return nil
	})
}
