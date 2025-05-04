package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
)

type cartItemService struct {
	cartItemRepository repository.CartItemRepository
	cartRepository     repository.CartRepository
}

type CartItemService interface {
	GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsWithQueryParamRequest) ([]model.CartItem, error)
	GetCartItemById(ctx context.Context, reqDTO *dto.GetCartItemByIdRequest) (*model.CartItem, error)
	GetCartItemsByCartId(ctx context.Context, reqDTO *dto.GetCartItemsByCartIdWithQueryParamRequest) ([]model.CartItem, error)
	CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) error
	UpdateCartItemById(ctx context.Context, reqDTO *dto.UpdateCartItemRequest) error
	DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteCartItemRequest) error
}

func NewCartItemService(cartItemRepository repository.CartItemRepository, cartRepository repository.CartRepository) CartItemService {
	return &cartItemService{
		cartItemRepository: cartItemRepository,
		cartRepository:     cartRepository,
	}
}

func (cartItemService *cartItemService) GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsWithQueryParamRequest) ([]model.CartItem, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	cartItemItems, err := cartItemService.cartItemRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return cartItemItems, nil
}

func (cartItemService *cartItemService) GetCartItemById(ctx context.Context, reqDTO *dto.GetCartItemByIdRequest) (*model.CartItem, error) {
	foundCartItem, err := cartItemService.cartItemRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, err
	}

	return foundCartItem, nil
}

func (cartItemService *cartItemService) GetCartItemsByCartId(ctx context.Context, reqDTO *dto.GetCartItemsByCartIdWithQueryParamRequest) ([]model.CartItem, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	cartItemItems, err := cartItemService.cartItemRepository.GetByCartId(ctx, reqDTO.CartId, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return cartItemItems, nil
}

func (cartItemService *cartItemService) CreateCartItem(ctx context.Context, reqDTO *dto.CreateCartItemRequest) error {
	foundCart, err := cartItemService.cartRepository.GetById(ctx, reqDTO.CartId)
	if err != nil {
		return fmt.Errorf("id of cart is not valid")
	}

	newCartItem := model.CartItem{
		CartId:    reqDTO.CartId,
		ProductId: reqDTO.Body.ProductId,
		Quantity:  1,
	}
	if err := cartItemService.cartItemRepository.Create(ctx, &newCartItem); err != nil {
		return err
	}

	if err := cartItemService.cartRepository.UpdateById(ctx, reqDTO.CartId, foundCart); err != nil {
		return err
	}

	return nil
}

func (cartItemService *cartItemService) UpdateCartItemById(ctx context.Context, reqDTO *dto.UpdateCartItemRequest) error {
	foundCart, err := cartItemService.cartRepository.GetById(ctx, reqDTO.CartId)
	if err != nil {
		return fmt.Errorf("id of cart is not valid")
	}

	foundCartItem, err := cartItemService.cartItemRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of cart item is not valid")
	}

	if reqDTO.Body.Quantity != nil {
		foundCartItem.Quantity = *reqDTO.Body.Quantity
	}

	if err := cartItemService.cartItemRepository.UpdateById(ctx, reqDTO.Id, foundCartItem); err != nil {
		return err
	}

	if err := cartItemService.cartRepository.UpdateById(ctx, reqDTO.CartId, foundCart); err != nil {
		return err
	}

	return nil
}

func (cartItemService *cartItemService) DeleteCartItemById(ctx context.Context, reqDTO *dto.DeleteCartItemRequest) error {
	foundCart, err := cartItemService.cartRepository.GetById(ctx, reqDTO.CartId)
	if err != nil {
		return fmt.Errorf("id of cart is not valid")
	}

	if _, err := cartItemService.cartItemRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of cart item is not valid")
	}

	if err := cartItemService.cartItemRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return err
	}

	if err := cartItemService.cartRepository.UpdateById(ctx, reqDTO.CartId, foundCart); err != nil {
		return err
	}

	return nil
}
