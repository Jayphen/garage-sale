package handlers

import (
	"fmt"

	"garagesale.jayphen.dev/internal/model"
	"garagesale.jayphen.dev/internal/utils"
	toast "garagesale.jayphen.dev/ui/components"
	components "garagesale.jayphen.dev/ui/components/cart"
	price "garagesale.jayphen.dev/ui/components/item"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func getSessionCart(session *sessions.Session) string {
	cartId, ok := session.Values["cart"].(string)
	if !ok {
		return ""
	}
	return cartId
}

func addToCart(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		itemId := c.FormValue("itemId")

		item, err := (&model.Item{}).FindItemById(e.App.Dao(), itemId)
		if err != nil {
			fmt.Println(err)
			return err
		}

		session := utils.GetSession(c.Request())
		cartId := getSessionCart(session)

		newCartId, err := model.AddToCart(e.App.Dao(), item.Id, cartId)
		if err != nil {
			fmt.Println(err)
			c.Response().Header().Set("HX-Reswap", "innerHTML")
			return utils.Render(c, 400, toast.Toast("That's already in your cart. You should probably hurry up and buy it"))
		}

		cartSize, err := model.GetCartSize(e.App.Dao(), newCartId)
		if err != nil {
			fmt.Println(err)
		}

		session.Values["cart"] = newCartId
		session.Values["cartSize"] = cartSize
		session.Save(c.Request(), c.Response())

		cartRecord, err := model.GetExistingCartRecord(e.App.Dao(), newCartId)
		if err != nil {
			return err
		}

		e.App.Dao().ExpandRecord(cartRecord, []string{"cartItems"}, nil)

		cart := model.Cart{
			Id:        cartId,
			Email:     "",
			CartItems: cartRecord.GetStringSlice("cartItems"),
		}

		expandedCart := model.NewExpandedCartFromCart(cart, cartRecord.ExpandedAll("cartItems"))

		return utils.Render(c, 200, components.CartSlideover(expandedCart))
	}
}

func getCartPreview(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		session := utils.GetSession(c.Request())
		cartId := getSessionCart(session)

		if cartId != "" {
			cartRecord, err := model.GetExistingCartRecord(e.App.Dao(), cartId)
			if err != nil {
				return err
			}

			e.App.Dao().ExpandRecord(cartRecord, []string{"cartItems"}, nil)

			cart := model.Cart{
				Id:        cartId,
				Email:     "",
				CartItems: cartRecord.GetStringSlice("cartItems"),
			}

			expandedCart := model.NewExpandedCartFromCart(cart, cartRecord.ExpandedAll("cartItems"))

			return utils.Render(c, 200, components.CartSlideover(expandedCart))

		}

		return nil
	}
}

func getCartTotal(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		session := utils.GetSession(c.Request())
		cartId := getSessionCart(session)

		if cartId != "" {
			cartRecord, err := model.GetExistingCartRecord(e.App.Dao(), cartId)
			if err != nil {
				return err
			}

			cart := model.Cart{
				Id:        cartRecord.GetString("id"),
				CartItems: cartRecord.GetStringSlice("cartItems"),
			}

			e.App.Dao().ExpandRecord(cartRecord, []string{"cartItems"}, nil)

			expandedCart := model.NewExpandedCartFromCart(cart, cartRecord.ExpandedAll("cartItems"))

			return utils.Render(c, 200, price.TotalPrice(expandedCart.TotalPrice))
		}

		return nil
	}
}

func removeFromCart(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		session := utils.GetSession(c.Request())
		cartId := getSessionCart(session)

		itemId := c.PathParam("item")

		if cartId != "" {
			cartRecord, err := model.GetExistingCartRecord(e.App.Dao(), cartId)
			if err != nil {
				return err
			}

			if err := model.RemoveFromCart(e.App.Dao(), cartRecord, itemId); err != nil {
				return err
			}

			cart := model.Cart{
				Id:        cartRecord.GetString("id"),
				CartItems: cartRecord.GetStringSlice("cartItems"),
			}

			e.App.Dao().ExpandRecord(cartRecord, []string{"cartItems"}, nil)
			expandedCart := model.NewExpandedCartFromCart(cart, cartRecord.ExpandedAll("cartItems"))

			return utils.Render(c, 200, components.CartContents(expandedCart))
		}

		return nil
	}
}

func RegisterCartHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/cart", addToCart(e))
		e.Router.DELETE("/cart/:item", removeFromCart(e))
		e.Router.GET("/cart-preview", getCartPreview(e))
		e.Router.GET("/cart-total", getCartTotal(e))
		return nil
	})
}
