package model

import (
	"errors"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

func AddToCart(dao *daos.Dao, id string, cartId string) (string, error) {
	item, err := dao.FindRecordById("items", id)
	if err != nil {
		return "", err
	}

	if cartId != "" {
		cartRecord, err := GetExistingCartRecord(dao, cartId)
		if err != nil {
			return "", err
		}
		return updateCartRecord(dao, cartRecord, item)
	} else {
		return createNewCartRecord(dao, item)
	}
}

func GetExistingCartRecord(dao *daos.Dao, cartId string) (*models.Record, error) {
	cartRecord, err := dao.FindRecordById("carts", cartId)
	if err != nil {
		return nil, err
	}
	return cartRecord, nil
}

func createNewCartRecord(dao *daos.Dao, item *models.Record) (string, error) {
	cartCollection, err := dao.FindCollectionByNameOrId("carts")
	if err != nil {
		return "", err
	}

	newCart := models.NewRecord(cartCollection)
	newCart.Set("cartItems", []string{item.Id})

	if err := dao.SaveRecord(newCart); err != nil {
		return "", err
	}

	return newCart.Id, nil
}

func updateCartRecord(dao *daos.Dao, cartRecord *models.Record, item *models.Record) (string, error) {
	// Check if the item is already in the cart
	existingCartItems := cartRecord.Get("cartItems")
	if existingCartItems != nil {
		items := existingCartItems.([]string)
		for _, existingItem := range items {
			if existingItem == item.Id {
				return "", errors.New("Item is already in the cart")
			}
		}
	}

	// Append the item ID to the existing cartItems array
	cartItems := append(existingCartItems.([]string), item.Id)

	// Update the cartItems field in the cart record
	cartRecord.Set("cartItems", cartItems)

	if err := dao.SaveRecord(cartRecord); err != nil {
		return "", err
	}

	return cartRecord.Id, nil
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
