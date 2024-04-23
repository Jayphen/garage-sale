package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/mail"

	"garagesale.jayphen.dev/internal/model"
	"garagesale.jayphen.dev/internal/utils"
	components "garagesale.jayphen.dev/ui/components/checkout"
	"garagesale.jayphen.dev/ui/pages"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
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

func sendConfirmationEmail(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		userEmail := c.FormValue("email")
		session := utils.GetSession(c.Request())

		if userEmail == "" {
			return echo.NewHTTPError(400, "Email is required")
		}

		// Retrieve cart ID from the session if it exists
		cartId, ok := session.Values["cart"].(string)
		if !ok {
			return echo.NewHTTPError(400, "It seems like your cart is empty")
		}

		if err := model.SetCartEmail(e.App.Dao(), cartId, userEmail); err != nil {
			return echo.NewHTTPError(500, "Could not save cart")
		}

		// todo: create a new entry in the db with a random token, used set to false and cartId set to the cartId

		tokenCollection, err := e.App.Dao().FindCollectionByNameOrId("confirmation_tokens")
		if err != nil {
			return echo.NewHTTPError(500, "Could not save cart")
		}

		token := make([]byte, 32/2) // Divide by 2 since hex encoding uses 2 characters for each byte
		if _, err := rand.Read(token); err != nil {
			return echo.NewHTTPError(500, "Could not save cart")
		}
		encodedToken := hex.EncodeToString(token)

		newTokenRecord := models.NewRecord(tokenCollection)
		newTokenRecord.Set("used", false)
		newTokenRecord.Set("cartId", cartId)
		newTokenRecord.Set("token", encodedToken)

		if err := e.App.Dao().SaveRecord(newTokenRecord); err != nil {
			return echo.NewHTTPError(500, "Could not save cart")
		}

		urlBase := e.App.Settings().Meta.AppUrl

		// todo: throttling, queue

		fmt.Println(urlBase)

		html := fmt.Sprintf("Hello my friend! Please click this link to verify your order: %s/confirm/%s", urlBase, encodedToken)

		fmt.Println(html)

		message := &mailer.Message{
			From: mail.Address{
				Address: e.App.Settings().Meta.SenderAddress,
				Name:    e.App.Settings().Meta.SenderName,
			},
			To:      []mail.Address{{Address: c.FormValue("email")}},
			Subject: "You bought stuff!",
			HTML:    html,
		}

		if err := e.App.NewMailClient().Send(message); err != nil {
			return echo.NewHTTPError(500, "Could not send email")
		}

		return utils.Render(c, 200, components.Thanks())
	}
}

func confirmToken(e *core.ServeEvent) func(echo.Context) error {
	return func(c echo.Context) error {
		token := c.PathParam("token")

		tokenRecord, err := e.App.Dao().FindFirstRecordByFilter("confirmation_tokens", "token={:token}", dbx.Params{"token": token})
		if err != nil {
			fmt.Errorf(err.Error())
			return echo.NewHTTPError(500, "Your link seems to be invalid")
		}

		cartRecord, err := model.GetExistingCartRecord(e.App.Dao(), tokenRecord.Get("cartId").(string))
		if err != nil {
			fmt.Errorf(err.Error())
			return echo.NewHTTPError(500, "Something went wrong")
		}

		cart := model.Cart{
			Id:        cartRecord.GetString("id"),
			CartItems: cartRecord.GetStringSlice("cartItems"),
		}
		e.App.Dao().ExpandRecord(cartRecord, []string{"cartItems"}, nil)
		expandedCart := model.NewExpandedCartFromCart(cart, cartRecord.ExpandedAll("cartItems"))

		return utils.Render(c, 200, pages.Confirm(expandedCart))
	}
}

func RegisterCheckoutHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/checkout", getCheckoutPage(e))
		e.Router.GET("/confirm/:token", confirmToken(e))
		e.Router.POST("/checkout", sendConfirmationEmail(e))
		return nil
	})
}
