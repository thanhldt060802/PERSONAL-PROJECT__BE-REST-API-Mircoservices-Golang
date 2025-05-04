package service

import (
	"context"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
)

type cartService struct {
	cartRepository repository.CartRepository
}

type CartService interface {
	GetCarts(ctx context.Context, reqDTO *dto.GetCartsWithQueryParamRequest) ([]model.Cart, error)
	GetCartByUserId(ctx context.Context, reqDTO *dto.GetCartByUserIdRequest) (*model.Cart, error)
}

func NewCartService(cartRepository repository.CartRepository) CartService {
	return &cartService{
		cartRepository: cartRepository,
	}
}

func (cartService *cartService) GetCarts(ctx context.Context, reqDTO *dto.GetCartsWithQueryParamRequest) ([]model.Cart, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	carts, err := cartService.cartRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return carts, nil
}

func (cartService *cartService) GetCartByUserId(ctx context.Context, reqDTO *dto.GetCartByUserIdRequest) (*model.Cart, error) {
	foundCart, err := cartService.cartRepository.GetByUserId(ctx, reqDTO.UserId)
	if err != nil {
		return nil, err
	}

	return foundCart, nil
}
