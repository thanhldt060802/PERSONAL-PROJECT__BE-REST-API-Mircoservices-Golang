package dto

import (
	"thanhldt060802/internal/model"
)

type CartItemView struct {
	Id        int64 `json:"id"`
	CartId    int64 `json:"cart_id"`
	ProductId int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func ToCartItemView(cartItem *model.CartItem) *CartItemView {
	return &CartItemView{
		Id:        cartItem.Id,
		CartId:    cartItem.CartId,
		ProductId: cartItem.ProductId,
		Quantity:  cartItem.Quantity,
	}
}

func ToListCartItemView(cartItems []model.CartItem) []CartItemView {
	cartItemViews := make([]CartItemView, len(cartItems))
	for i, cartItem := range cartItems {
		cartItemViews[i] = *ToCartItemView(&cartItem)
	}
	return cartItemViews
}
