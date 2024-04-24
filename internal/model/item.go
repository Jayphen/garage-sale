package model

import (
	"fmt"

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

	return items, nil
}

func (item *Item) FindItemById(dao *daos.Dao) (*Item, error) {
	err := ItemQuery(dao).
		AndWhere(dbx.HashExp{"id": item.Id}).
		One(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (item *Item) SetItemStatus(dao *daos.Dao, status ItemStatus) error {
	record, err := dao.FindRecordById("items", item.Id)
	if err != nil {
		return err
	}

	record.Set("status", string(status))

	if string(status) == "sold" {
		record.Set("sellPrice", record.Get("price"))
	}

	if err := dao.SaveRecord(record); err != nil {
		return err
	}

	return nil
}
