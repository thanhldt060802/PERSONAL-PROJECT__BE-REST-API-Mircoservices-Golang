package dto

import (
	"strings"
	"thanhldt060802/internal/model"
	"time"

	"github.com/shopspring/decimal"
)

// Struct to parse sorting field in query
// ################################################################################
type SortField struct {
	Field     string
	Direction string
}

func ParseSortBy(sortBy string) []SortField {
	var sortFields []SortField

	items := strings.Split(sortBy, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		var field string
		var direction string

		if strings.Contains(item, ":") {
			parts := strings.SplitN(item, ":", 2)
			field = parts[0]
			if strings.ToLower(parts[1]) == "desc" {
				direction = "DESC"
			} else {
				direction = "ASC"
			}
		} else {
			field = item
			direction = "ASC"
		}

		sortFields = append(sortFields, SortField{
			Field:     field,
			Direction: direction,
		})
	}

	return sortFields
}

// ################################################################################

// DTO for data responding
// ################################################################################
type ProductDTO struct {
	Id                 int64           `json:"id"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Price              decimal.Decimal `json:"price"`
	DiscountPercentage int32           `json:"discount_percentage"`
	Stock              int32           `json:"stock"`
	ImageURL           string          `json:"image_url"`
	CategoryId         int64           `json:"category_id"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

func ToProductDTO(product *model.Product) *ProductDTO {
	return &ProductDTO{
		Id:                 product.Id,
		Name:               product.Name,
		Description:        product.Description,
		Price:              product.Price,
		DiscountPercentage: product.DiscountPercentage,
		Stock:              product.Stock,
		ImageURL:           product.ImageURL,
		CategoryId:         product.CategoryId,
		CreatedAt:          product.CreatedAt,
		UpdatedAt:          product.UpdatedAt,
	}
}

func ToProductDTOs(products []model.Product) []ProductDTO {
	productDTOs := make([]ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = *ToProductDTO(&product)
	}
	return productDTOs
}

// ################################################################################

// Request
// ################################################################################
type GetProductsRequestQueryParam struct {
	Offset int    `query:"offset" default:"0" minimum:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"name:asc" example:"name,price:desc" doc:"Sort by one or more fields separated by commas. For example: sort_by=name,price:desc will sort by name in ascending order, then by price in descending order."`
}

type GetProductByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of product will be gotten."`
}

type CreateProductRequest struct {
	Body struct {
		Name               string          `json:"name" required:"true" minLength:"1" doc:"Name of product."`
		Description        string          `json:"description" required:"true" minLength:"1" doc:"Description of product."`
		Price              decimal.Decimal `json:"price" required:"true" minimum:"0" doc:"Price of product."`
		DiscountPercentage int32           `json:"discount_percentage" required:"true" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              int32           `json:"stock" required:"true" minimun:"0" doc:"Stock of product."`
		ImageURL           string          `json:"image_url" required:"true" minLength:"1" doc:"Image URL of product."`
		CategoryId         int64           `json:"category_id" required:"true" minimum:"1" doc:"Category id of product."`
	}
}

type UpdateProductRequest struct {
	Id   int64 `path:"id" required:"true"`
	Body struct {
		Name               *string          `json:"name" required:"true" minLength:"1" doc:"Name of product."`
		Description        *string          `json:"description" required:"true" minLength:"1" doc:"Description of product."`
		Price              *decimal.Decimal `json:"price" required:"true" minimum:"0" doc:"Price of product."`
		DiscountPercentage *int32           `json:"discount_percentage" required:"true" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              *int32           `json:"stock" required:"true" minimun:"0" doc:"Stock of product."`
		ImageURL           *string          `json:"image_url" required:"true" minLength:"1" doc:"Image URL of product."`
		CategoryId         *int64           `json:"category_id" required:"true" minimum:"1" doc:"Category id of product."`
	}
}

type DeleteProductRequest struct {
	Id int64 `path:"id" required:"true"`
}

// ################################################################################
