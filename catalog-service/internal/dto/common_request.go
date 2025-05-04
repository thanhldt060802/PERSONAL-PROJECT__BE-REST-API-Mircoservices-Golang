package dto

import "github.com/shopspring/decimal"

// Only category request
// ################################################################################

type GetCategoriesWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetCategoryByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of category will be gotten."`
}

type CreateCategoryRequest struct {
	Body struct {
		Name        string `json:"name" required:"true" minLength:"1" doc:"Name of category."`
		Description string `json:"description" required:"true" minLength:"1" doc:"Description of category."`
	}
}

type UpdateCategoryRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of category will be updated."`
	Body struct {
		Name        *string `json:"name,omitempty" minLength:"1" doc:"Name of category."`
		Description *string `json:"description,omitempty" minLength:"1" doc:"Description of category."`
	}
}

type DeleteCategoryRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of category will be deleted."`
}

// ################################################################################

// Only product request
// ################################################################################

type GetProductsWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetProductByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of product will be gotten."`
}

type GetProductsByCategoryIdWithQueryParamRequest struct {
	CategoryId int64  `path:"category_id" required:"true" doc:"Id of category will be filtered."`
	Offset     int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit      int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy     string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type CreateProductRequest struct {
	Body struct {
		Name               string          `json:"name" required:"true" minLength:"1" doc:"Name of product."`
		Description        string          `json:"description" required:"true" minLength:"1" doc:"Description of product."`
		Price              decimal.Decimal `json:"price" required:"true" doc:"Price of product."`
		DiscountPercentage int32           `json:"discount_percentage" required:"true" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              int32           `json:"stock" required:"true" minimun:"0" doc:"Stock of product."`
		ImageURL           string          `json:"image_url" required:"true" minLength:"1" doc:"Image URL of product."`
		CategoryId         int64           `json:"category_id" required:"true" minimum:"1" doc:"Category id of product."`
	}
}

type UpdateProductRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of product will be updated."`
	Body struct {
		Name               *string          `json:"name,omitempty" minLength:"1" doc:"Name of product."`
		Description        *string          `json:"description,omitempty" minLength:"1" doc:"Description of product."`
		Price              *decimal.Decimal `json:"price,omitempty" doc:"Price of product."`
		DiscountPercentage *int32           `json:"discount_percentage,omitempty" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              *int32           `json:"stock,omitempty" minimun:"0" doc:"Stock of product."`
		ImageURL           *string          `json:"image_url,omitempty" minLength:"1" doc:"Image URL of product."`
		CategoryId         *int64           `json:"category_id,omitempty" minimum:"1" doc:"Category id of product."`
	}
}

type DeleteProductRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of product will be deleted."`
}

type SearchProductsRequest struct {
	Query string `query:"query" required:"true" doc:"Search query for products."`
}

// ################################################################################
