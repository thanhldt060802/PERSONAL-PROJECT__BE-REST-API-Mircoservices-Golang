package dto

import (
	"thanhldt060802/internal/model"
	"time"
)

type CartView struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCartView(cart *model.Cart) *CartView {
	return &CartView{
		Id:        cart.Id,
		UserId:    cart.UserId,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}

func ToListCartView(carts []model.Cart) []CartView {
	cartViews := make([]CartView, len(carts))
	for i, cart := range carts {
		cartViews[i] = *ToCartView(&cart)
	}
	return cartViews
}
