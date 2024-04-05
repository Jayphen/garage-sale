package model

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Bid struct {
	models.BaseModel
	BidderEmail string `db:"bidder_email" json:"bidderEmail"`
	ItemId      string `db:"item_id" json:"itemId"`
	Amount      int    `db:"amount" json:"amount"`
}

var _ models.Model = (*Bid)(nil)

func (m *Bid) TableName() string {
	return "bids"
}

func (bid *Bid) CreateBid(dao *daos.Dao) error {
	// should validate...
	return dao.Save(bid)
}
