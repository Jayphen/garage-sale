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
}

var _ models.Model = (*Item)(nil)

func (m *Item) TableName() string {
	return "items"
}

func ItemQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Item{})
}

func (item *Item) GetItems(dao *daos.Dao) ([]*Item, error) {
	var c []*Item
	err := ItemQuery(dao).OrderBy("title asc").All(&c)
	return c, err
}
