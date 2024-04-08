package model

import (
	"errors"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func AddToCart(dao *daos.Dao, id string) (string, error) {
	item, err := dao.FindRecordById("items", id)
	if err != nil {
		return "", err
	}

	cartCollection, err := dao.FindCollectionByNameOrId("carts")
	if err != nil {
		return "", err
	}

	cart := models.NewRecord(cartCollection)

	cart.Set("cartItems", []string{item.Id})

	if err := dao.SaveRecord(cart); err != nil {
		return "", err
	}

	return cart.Id, nil
}

func GetCartSize(dao *daos.Dao, id string) (int, error) {
	cart, err := dao.FindRecordById("carts", id)
	if err != nil {
		return 0, err
	}

	// Check if the cartItems value is a slice of strings
	items, ok := cart.Get("cartItems").([]string)
	if !ok {
		return 0, errors.New("cartItems is not a slice of strings")
	}

	// Get the length of the items array
	cartSize := len(items)

	return cartSize, nil
}
