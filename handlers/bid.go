package handlers

import (
	"garagesale.jayphen.dev/assets/templ/pages"
	"garagesale.jayphen.dev/model"
	"garagesale.jayphen.dev/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/models"
)

func RegisterBidHandlers(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("bid", func(c echo.Context) error {
			postData := struct {
				BidderEmail string `json:"bidder_email" form:"bidder_email" validate:"required"`
				ItemId      string `json:"item_id" form:"item_id" validate:"required,uuid"`
				Amount      string `json:"amount" form:"amount" validate:"required,gt=0"`
			}{}

			if err := c.Bind(&postData); err != nil {
				return apis.NewBadRequestError("Failed to parse form data", err)
			}

			collection, err := app.Dao().FindCollectionByNameOrId("bids")
			if err != nil {
				return apis.NewNotFoundError("Failed to find bids collection", err)
			}

			record := models.NewRecord(collection)

			form := forms.NewRecordUpsert(app, record)

			form.LoadData(map[string]any{
				"bidder_email": postData.BidderEmail,
				"item_id":      postData.ItemId,
				"amount":       postData.Amount,
			})

			if err := form.Submit(); err != nil {
				// maybe check for valid item here

				return apis.NewBadRequestError("Failed to submit bid", err)
			}

			item, err := (&model.Item{}).FindItemById(app.Dao(), postData.ItemId)
			if err != nil {
				return err
			}

			{
				var bids []*model.Bid

				err := app.Dao().ModelQuery(&model.Bid{}).Where(dbx.HashExp{"item_id": item.Id}).All(&bids)
				if err != nil {
					return err
				}
				item.Bids = bids
			}

			return utils.Render(c, 201, pages.Price(item))
		})

		return nil
	})
}
