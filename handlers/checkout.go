package handlers

import (
	"garagesale.jayphen.dev/model"
	"garagesale.jayphen.dev/ui/pages"
	"garagesale.jayphen.dev/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func getCheckoutPage(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		session := utils.GetSession(c.Request())

		// Retrieve cart ID from the session if it exists
		cartId, ok := session.Values["cart"].(string)
		if !ok {
			cartId = ""
		}

		if cartId != "" {
			cartRecord, err := model.GetExistingCartRecord(e.App.Dao(), cartId)
			if err != nil {
				return err
			}

			e.App.Dao().ExpandRecord(cartRecord, []string{"cartItems"}, nil)

			cart := model.Cart{
				Id:        cartRecord.GetString("id"),
				CartItems: cartRecord.GetStringSlice("cartItems"),
			}
			expandedCart := model.NewExpandedCartFromCart(cart, cartRecord.ExpandedAll("cartItems"))

			// return the checkout ui
			return utils.Render(c, 200, pages.Checkout(&expandedCart))

		}

		// no cart, return empty cart ui
		return nil
	}
}

func RegisterCheckoutHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/checkout", getCheckoutPage(e))
		return nil
	})
}
