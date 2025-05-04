package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type CartHandler struct {
	cartService    service.CartService
	authMiddleware *middleware.AuthMiddleware
}

func NewCartHandler(api huma.API, cartService service.CartService, authMiddleware *middleware.AuthMiddleware) *CartHandler {
	cartHandler := &CartHandler{
		cartService:    cartService,
		authMiddleware: authMiddleware,
	}

	// Get carts
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/carts",
		Summary:     "/carts",
		Description: "Get carts.",
		Tags:        []string{"Cart"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, cartHandler.GetCarts)

	// Get cart by user id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/carts/user-id/{user_id}",
		Summary:     "/carts/user-id/{user_id}",
		Description: "Get cart by user id.",
		Tags:        []string{"Cart"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, cartHandler.GetCartByUserId)

	// Get cart using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-cart",
		Summary:     "/my-cart",
		Description: "Get cart using account.",
		Tags:        []string{"Cart"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, cartHandler.GetCartUsingAccount)

	return cartHandler
}

func (cartHandler *CartHandler) GetCarts(ctx context.Context, reqDTO *dto.GetCartsWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.CartView], error) {
	carts, err := cartHandler.cartService.GetCarts(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get carts failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListCartView(carts)
	res := &dto.PaginationBodyResponseList[dto.CartView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get carts successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (cartHandler *CartHandler) GetCartByUserId(ctx context.Context, reqDTO *dto.GetCartByUserIdRequest) (*dto.BodyResponse[dto.CartView], error) {
	foundCart, err := cartHandler.cartService.GetCartByUserId(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get cart by user id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToCartView(foundCart)
	res := &dto.BodyResponse[dto.CartView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart by user id successful"
	res.Body.Data = *data
	return res, nil
}

func (cartHandler *CartHandler) GetCartUsingAccount(ctx context.Context, reqDTO *struct{}) (*dto.BodyResponse[dto.CartView], error) {
	userId := ctx.Value("cart_id").(int64)

	convertReqDTO := &dto.GetCartByUserIdRequest{UserId: userId}

	foundCart, err := cartHandler.cartService.GetCartByUserId(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get cart using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToCartView(foundCart)
	res := &dto.BodyResponse[dto.CartView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart using account successful"
	res.Body.Data = *data
	return res, nil
}
