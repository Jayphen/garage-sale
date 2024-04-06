package model

import (
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Bid struct {
	models.BaseModel
	Id          string `db:"id" json:"id"`
	BidderEmail string `db:"bidder_email" json:"bidder_email"`
	ItemId      string `db:"item_id" json:"item_id"`
	Amount      string `db:"amount" json:"amount"`
}

var _ models.Model = (*Bid)(nil)

func (m *Bid) TableName() string {
	return "bids"
}

func (bid *Bid) CreateBid(dao *daos.Dao) error {
	// should validate...
	return dao.Save(bid)
}
