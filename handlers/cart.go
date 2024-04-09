package handlers

import (
	"fmt"

	"garagesale.jayphen.dev/model"
	toast "garagesale.jayphen.dev/ui/components"
	components "garagesale.jayphen.dev/ui/components/cart"
	price "garagesale.jayphen.dev/ui/components/item"
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
			return utils.Render(c, 400, toast.Toast("That's already in your cart. You should probably hurry up and buy it"))
		}

		cartSize, err := model.GetCartSize(e.App.Dao(), newCartId)
		if err != nil {
			fmt.Println(err)
		}

		session.Values["cart"] = newCartId
		session.Values["cartSize"] = cartSize
		session.Save(c.Request(), c.Response())

		return utils.Render(c, 200, components.Indicator(cartSize))
	}
}

func getCartPreview(e *core.ServeEvent) func(echo.Context) error {
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

			var cartItems []components.CartItem
			for _, item := range cartRecord.ExpandedAll("cartItems") {
				cartItem := components.CartItem{
					Title:  item.GetString("title"),
					Price:  item.GetInt("price"),
					Id:     item.GetString("id"),
					Images: item.GetStringSlice("images"),
				}
				cartItems = append(cartItems, cartItem)
			}

			return utils.Render(c, 200, components.CartSlideover(cartItems))

		}

		return nil
	}
}

func getCartTotal(e *core.ServeEvent) func(echo.Context) error {
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

			var itemPrices []price.ItemPrice
			for _, item := range cartRecord.ExpandedAll("cartItems") {
				itemPrice := price.ItemPrice{
					Price: item.GetInt("price"),
					Id:    item.GetString("id"),
				}
				itemPrices = append(itemPrices, itemPrice)
			}

			return utils.Render(c, 200, price.TotalPrice(itemPrices))

		}

		return nil
	}
}

func RegisterCartHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/cart", addToCart(e))
		e.Router.GET("/cart-preview", getCartPreview(e))
		e.Router.GET("/cart-total", getCartTotal(e))
		return nil
	})
}
