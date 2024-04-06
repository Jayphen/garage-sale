package model

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Item struct {
	models.BaseModel
	Id          string                  `db:"id" json:"id"`
	Title       string                  `db:"title" json:"title"`
	Description string                  `db:"description" json:"description"`
	Price       int                     `db:"price" json:"price"`
	Images      types.JsonArray[string] `db:"images" json:"images"`
	Bids        []*Bid                  `db:"bids" json:"bids"`
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

	err := ItemQuery(dao).OrderBy("title asc").All(&items)
	if err != nil {
		return nil, err
	}

	// n+1 but it is sqlite so no worries
	for _, item := range items {
		var bids []*Bid

		err := dao.ModelQuery(&Bid{}).Where(dbx.HashExp{"item_id": item.Id}).All(&bids)
		if err != nil {
			return nil, err
		}

		item.Bids = bids
	}

	return items, nil
}

func (item *Item) FindItemById(dao *daos.Dao, id string) (*Item, error) {
	err := ItemQuery(dao).
		AndWhere(dbx.HashExp{"id": id}).
		Limit(1).
		One(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}
