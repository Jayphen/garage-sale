package model

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type Cart struct {
	models.BaseModel
	Id        string   `db:"id" json:"id"`
	CartItems []string `db:"cartItems" json:"cartItems"`
	Email     string   `db:"email" json:"email"`
}

type ExpandedCart struct {
	Cart
	CartItems  []Item `db:"cartItems" json:"cartItems"`
	TotalPrice string
	FinalPrice string
	CartSize   int
}

func NewExpandedCartFromCart(cart Cart, items []*models.Record) ExpandedCart {
	cartSize := len(cart.CartItems)

	expandedCart := &ExpandedCart{
		Cart:      cart,
		CartItems: make([]Item, cartSize),
	}

	for i, item := range items {
		cartItem := Item{
			Title:     item.GetString("title"),
			Price:     item.GetInt("price"),
			SellPrice: item.GetInt("sellPrice"),
			Id:        item.GetString("id"),
			Images:    item.GetStringSlice("images"),
			MinPrice:  item.GetInt("minPrice"),
			MaxPrice:  item.GetInt("maxPrice"),
		}

		cartItem.Price = cartItem.CalculateCurrentPrice()

		expandedCart.CartItems[i] = cartItem
	}

	expandedCart.getTotalPrice()
	expandedCart.CartSize = cartSize

	return *expandedCart
}

func (c *ExpandedCart) getTotalPrice() {
	total := 0
	finalTotal := 0
	for _, item := range c.CartItems {
		total += item.Price
		finalTotal += item.SellPrice
	}

	fmt.Println(finalTotal)
	c.TotalPrice = strconv.FormatFloat(float64(total)/100.0, 'f', 2, 64)
	c.FinalPrice = strconv.FormatFloat(float64(finalTotal)/100.0, 'f', 2, 64)
}

var _ models.Model = (*Cart)(nil)

func (m *Cart) TableName() string {
	return "carts"
}

type Price struct {
	Price int
}

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

func RemoveFromCart(dao *daos.Dao, cartRecord *models.Record, itemId string) error {
	// Check if the item is already in the cart
	existingCartItems, ok := cartRecord.Get("cartItems").([]string)
	if !ok || existingCartItems == nil {
		return errors.New("Cart is empty")
	}

	var updatedCartItems []string
	var itemFound bool

	for i, item := range existingCartItems {
		if item == itemId {
			updatedCartItems = append(existingCartItems[:i], existingCartItems[i+1:]...)
			cartRecord.Set("cartItems", updatedCartItems)
			itemFound = true
			break
		}
	}

	if !itemFound {
		return errors.New("Item not found in cart")
	}

	if err := dao.SaveRecord(cartRecord); err != nil {
		return err
	}

	return nil
}

func GetExistingCartRecord(dao *daos.Dao, cartId string) (*models.Record, error) {
	cartRecord, err := dao.FindRecordById("carts", cartId)
	if err != nil {
		return nil, err
	}
	return cartRecord, nil
}

func SetCartEmail(dao *daos.Dao, cartId string, email string) error {
	cartRecord, err := dao.FindRecordById("carts", cartId)
	if err != nil {
		return err
	}

	cartRecord.Set("email", email)

	if err := dao.SaveRecord(cartRecord); err != nil {
		return err
	}

	return nil
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
