package model

import (
	"fmt"
	"math"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ItemStatus string

const (
	StatusSoon      ItemStatus = "soon"
	StatusAvailable ItemStatus = "available"
	StatusFrozen    ItemStatus = "frozen"
	StatusSold      ItemStatus = "sold"
)

func (is *ItemStatus) ParseFormValue(value string) error {
	switch value {
	case "soon":
		*is = StatusSoon
	case "available":
		*is = StatusAvailable
	case "frozen":
		*is = StatusFrozen
	case "sold":
		*is = StatusSold
	default:
		return fmt.Errorf("Invalid value for ItemStatus: %s", value)
	}

	return nil
}

type Item struct {
	models.BaseModel
	Id               string                  `db:"id" json:"id"`
	Title            string                  `db:"title" json:"title"`
	Description      string                  `db:"description" json:"description"`
	ShortDescription string                  `db:"shortDesc" json:"shortDesc"`
	Price            int                     `db:"price" json:"price"`
	MinPrice         int                     `db:"minPrice" json:"minPrice"`
	MaxPrice         int                     `db:"maxPrice" json:"maxPrice"`
	SellPrice        int                     `db:"sellPrice" json:"sellPrice"`
	Images           types.JsonArray[string] `db:"images" json:"images"`
	Status           ItemStatus              `db:"status" json:"status"`
}

var _ models.Model = (*Item)(nil)

func (m *Item) TableName() string {
	return "items"
}

func ItemQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Item{})
}

func (item *Item) GetItems(dao *daos.Dao) ([]*Item, error) {
	var items []*Item

	err := ItemQuery(dao).OrderBy("created desc").All(&items)
	if err != nil {
		return nil, err
	}

	for _, i := range items {
		price := i.CalculateCurrentPrice()
		i.Price = price
	}

	return items, nil
}

func (item *Item) FindItemById(dao *daos.Dao) (*Item, error) {
	err := ItemQuery(dao).
		AndWhere(dbx.HashExp{"id": item.Id}).
		One(item)
	if err != nil {
		return nil, err
	}

	price := item.CalculateCurrentPrice()
	item.Price = price

	return item, nil
}

func (item *Item) SetItemStatus(dao *daos.Dao, status ItemStatus) error {
	record, err := dao.FindRecordById("items", item.Id)
	if err != nil {
		return err
	}

	price := item.CalculateCurrentPrice()
	item.Price = price

	record.Set("status", string(status))

	if string(status) == "sold" {
		record.Set("sellPrice", price)
	}

	if err := dao.SaveRecord(record); err != nil {
		return err
	}

	return nil
}

const (
	OperationalStartHour = 7
	OperationalEndHour   = 14
	timeFormat           = "2006-01-02 15:04:05.999Z"
)

// CalculateCurrentPrice calculates the current price based on the maximum, minimum prices,
// and the added time of the product.
func (item *Item) CalculateCurrentPrice() int {
	currentTime := time.Now()
	todayStartOperational := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), OperationalStartHour, 0, 0, 0, currentTime.Location())

	// Parse the addedTimeStr
	addedTimestamp, _ := time.Parse(timeFormat, item.Created.String())

	// Check if the current time is within the operational hours
	if currentTime.Hour() < OperationalStartHour || currentTime.Hour() >= OperationalEndHour {
		return item.MaxPrice // Out of operational hours
	}

	if addedTimestamp.Before(todayStartOperational) {
		addedTimestamp = todayStartOperational
	}

	// Calculate the total seconds for price adjustment
	totalOperationalSeconds := float64((OperationalEndHour - OperationalStartHour) * 3600)
	elapsedSeconds := currentTime.Sub(addedTimestamp).Seconds()

	if elapsedSeconds > totalOperationalSeconds {
		return item.MinPrice // Maximum time elapsed, return minimum price
	}

	// Linear price decrease calculation
	decreasePerSecond := (float64(item.MaxPrice) - float64(item.MinPrice)) / totalOperationalSeconds
	currentPrice := float64(item.MaxPrice) - (decreasePerSecond * elapsedSeconds)

	if currentPrice < float64(item.MinPrice) {
		currentPrice = float64(item.MinPrice)
	}

	return int(math.Round(currentPrice)) // Assuming prices are in cents and should be rounded to nearest whole number
}
