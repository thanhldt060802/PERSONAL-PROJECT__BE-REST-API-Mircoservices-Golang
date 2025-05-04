package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type CartItemHandler struct {
	cartItemService service.CartItemService
	authMiddleware  *middleware.AuthMiddleware
}

func NewCartItemHandler(api huma.API, cartItemService service.CartItemService, authMiddleware *middleware.AuthMiddleware) *CartItemHandler {
	cartItemHandler := &CartItemHandler{
		cartItemService: cartItemService,
		authMiddleware:  authMiddleware,
	}

	// Get cart items
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/cart-items",
		Summary:     "/cart-items",
		Description: "Get cart items.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, cartItemHandler.GetCartItems)

	// Get cart item by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/cart-items/id/{id}",
		Summary:     "/cart-items/id/{id}",
		Description: "Get cart item by id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, cartItemHandler.GetCartItemById)

	// Get cart items by cart id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/cart-items/cart-id/{cart_id}",
		Summary:     "/cart-items/cart-id/{cart_id}",
		Description: "Get cart items by cart id.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, cartItemHandler.GetCartItemsByCartId)

	// Get cart items using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-cart-items",
		Summary:     "/my-cart-items",
		Description: "Get cart items using account.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, cartItemHandler.GetCartItemsUsingAccount)

	// Get cart item by id using account
	// huma.Register(api, huma.Operation{
	// 	Method:      http.MethodGet,
	// 	Path:        "/my-cart-items/id/{id}",
	// 	Summary:     "/my-cart-items/id/{id}",
	// 	Description: "Get cart item by id using account.",
	// 	Tags:        []string{"Cart Item"},
	// 	Middlewares: huma.Middlewares{authMiddleware.Authentication},
	// }, cartItemHandler.GetCartItemByIdUsingAccount)

	// Create cart item using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/my-cart-items",
		Summary:     "/my-cart-items",
		Description: "Create cart item using account.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, cartItemHandler.CreateCartItemUsingAccount)

	// Update cart item by id using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/my-cart-items/id/{id}",
		Summary:     "/my-cart-items/id/{id}",
		Description: "Update cart item by id using account.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, cartItemHandler.UpdateCartItemUsingAccount)

	// Delete cart item by id using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/my-cart-items/id/{id}",
		Summary:     "/my-cart-items/id/{id}",
		Description: "Delete cart item by id using account.",
		Tags:        []string{"Cart Item"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, cartItemHandler.DeleteCartItemUsingAccount)

	return cartItemHandler
}

func (cartItemHandler *CartItemHandler) GetCartItems(ctx context.Context, reqDTO *dto.GetCartItemsWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.CartItemView], error) {
	cartItems, err := cartItemHandler.cartItemService.GetCartItems(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get cart items failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListCartItemView(cartItems)
	res := &dto.PaginationBodyResponseList[dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart items successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (cartItemHandler *CartItemHandler) GetCartItemById(ctx context.Context, reqDTO *dto.GetCartItemByIdRequest) (*dto.BodyResponse[dto.CartItemView], error) {
	foundCartItem, err := cartItemHandler.cartItemService.GetCartItemById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get cart item by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToCartItemView(foundCartItem)
	res := &dto.BodyResponse[dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart item by id successful"
	res.Body.Data = *data
	return res, nil
}

func (cartItemHandler *CartItemHandler) GetCartItemsByCartId(ctx context.Context, reqDTO *dto.GetCartItemsByCartIdWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.CartItemView], error) {
	cartItems, err := cartItemHandler.cartItemService.GetCartItemsByCartId(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get cart items by cart id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListCartItemView(cartItems)
	res := &dto.PaginationBodyResponseList[dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart items by cart id successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (cartItemHandler *CartItemHandler) GetCartItemsUsingAccount(ctx context.Context, reqDTO *dto.GetCartItemsUsingAccountWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.CartItemView], error) {
	cartId := ctx.Value("cart_id").(int64)

	convertReqDTO := &dto.GetCartItemsByCartIdWithQueryParamRequest{
		CartId: cartId,
		Offset: reqDTO.Offset,
		Limit:  reqDTO.Limit,
		SortBy: reqDTO.SortBy,
	}

	cartItems, err := cartItemHandler.cartItemService.GetCartItemsByCartId(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get cart items using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListCartItemView(cartItems)
	res := &dto.PaginationBodyResponseList[dto.CartItemView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get cart items using account successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

// func (cartItemHandler *CartItemHandler) GetCartItemByIdUsingAccount(ctx context.Context, reqDTO *dto.GetCartItemByIdUsingAccountRequest) (*dto.BodyResponse[dto.CartItemView], error) {
// 	cartId := ctx.Value("cart_id").(int64)

// 	convertReqDTO := &dto.GetCartItemByIdRequest{Id: reqDTO.Id}

// 	foundCartItem, err := cartItemHandler.cartItemService.GetCartItemById(ctx, convertReqDTO)
// 	if err != nil {
// 		res := &dto.ErrorResponse{}
// 		res.Status = http.StatusBadRequest
// 		res.Code = "ERR_BAD_REQUEST"
// 		res.Message = "Get cart item by id using account failed"
// 		res.Details = []string{err.Error()}
// 		return nil, res
// 	} else if foundCartItem.CartId != cartId {
// 		res := &dto.ErrorResponse{}
// 		res.Status = http.StatusBadRequest
// 		res.Code = "ERR_BAD_REQUEST"
// 		res.Message = "Get cart item by id using account failed"
// 		res.Details = []string{"id of cart item is not valid"}
// 		return nil, res
// 	}

// 	data := dto.ToCartItemView(foundCartItem)
// 	res := &dto.BodyResponse[dto.CartItemView]{}
// 	res.Body.Code = "OK"
// 	res.Body.Message = "Get cart item by id using account successful"
// 	res.Body.Data = *data
// 	return res, nil
// }

func (cartItemHandler *CartItemHandler) CreateCartItemUsingAccount(ctx context.Context, reqDTO *dto.CreateCartItemUsingAccountRequest) (*dto.SuccessResponse, error) {
	cartId := ctx.Value("cart_id").(int64)

	convertReqDTO := &dto.CreateCartItemRequest{CartId: cartId}
	convertReqDTO.Body.ProductId = reqDTO.Body.ProductId

	if err := cartItemHandler.cartItemService.CreateCartItem(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create cart item using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create cart item using account successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) UpdateCartItemUsingAccount(ctx context.Context, reqDTO *dto.UpdateCartItemUsingAccountRequest) (*dto.SuccessResponse, error) {
	cartId := ctx.Value("cart_id").(int64)

	foundCartItem, err := cartItemHandler.cartItemService.GetCartItemById(ctx, &dto.GetCartItemByIdRequest{Id: reqDTO.Id})
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update cart item using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	} else if foundCartItem.CartId != cartId {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update cart item using account failed"
		res.Details = []string{"id of cart item is not valid"}
		return nil, res
	}

	convertReqDTO := &dto.UpdateCartItemRequest{
		CartId: cartId,
		Id:     reqDTO.Id,
	}
	convertReqDTO.Body.Quantity = reqDTO.Body.Quantity

	if err := cartItemHandler.cartItemService.UpdateCartItemById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update cart item using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update cart item using account successful"
	return res, nil
}

func (cartItemHandler *CartItemHandler) DeleteCartItemUsingAccount(ctx context.Context, reqDTO *dto.DeleteCartItemUsingAccountRequest) (*dto.SuccessResponse, error) {
	cartId := ctx.Value("cart_id").(int64)

	foundCartItem, err := cartItemHandler.cartItemService.GetCartItemById(ctx, &dto.GetCartItemByIdRequest{Id: reqDTO.Id})
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete cart item using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	} else if foundCartItem.CartId != cartId {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete cart item using account failed"
		res.Details = []string{"id of cart item is not valid"}
		return nil, res
	}

	convertReqDTO := &dto.DeleteCartItemRequest{
		CartId: cartId,
		Id:     reqDTO.Id,
	}

	if err := cartItemHandler.cartItemService.DeleteCartItemById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete cart item using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete cart item using account successful"
	return res, nil
}
