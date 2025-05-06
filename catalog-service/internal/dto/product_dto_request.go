package dto

type GetProductsRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:desc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. Format: \"field:asc/desc\" (default is asc if not declare after commas)"`
}

type GetProductByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of product."`
}

type GetProductsByCategoryIdRequest struct {
	CategoryId int64  `path:"category_id" required:"true" doc:"Id of category."`
	Offset     int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit      int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy     string `query:"sort_by" default:"id:desc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. Format: \"field:asc/desc\" (default is asc if not declare after commas)"`
}

type CreateProductRequest struct {
	Body struct {
		Name               string `json:"name" required:"true" minLength:"1" doc:"Name of product."`
		Description        string `json:"description" required:"true" minLength:"1" doc:"Description of product."`
		Sex                string `json:"sex" required:"true" minLength:"1" enum:"MALE,FEMALE,UNISEX" doc:"Sex of product."`
		Price              int64  `json:"price" required:"true" minimum:"0" doc:"Price of product."`
		DiscountPercentage int32  `json:"discount_percentage" required:"true" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              int32  `json:"stock" required:"true" minimun:"0" doc:"Stock of product."`
		ImageURL           string `json:"image_url" required:"true" minLength:"1" doc:"Image URL of product."`
		CategoryId         int64  `json:"category_id" required:"true" minimum:"1" doc:"Category id of product."`
	}
}

type UpdateProductByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of product."`
	Body struct {
		Name               *string `json:"name,omitempty" minLength:"1" doc:"Name of product."`
		Description        *string `json:"description,omitempty" minLength:"1" doc:"Description of product."`
		Sex                *string `json:"sex,omitempty" minLength:"1" enum:"MALE,FEMALE,UNISEX" doc:"Sex of product."`
		Price              *int64  `json:"price,omitempty" minimum:"0" doc:"Price of product."`
		DiscountPercentage *int32  `json:"discount_percentage,omitempty" minimum:"0" maximum:"100" doc:"Discount percentage of product."`
		Stock              *int32  `json:"stock,omitempty" minimun:"0" doc:"Stock of product."`
		ImageURL           *string `json:"image_url,omitempty" minLength:"1" doc:"Image URL of product."`
		CategoryId         *int64  `json:"category_id,omitempty" minimum:"1" doc:"Category id of product."`
	}
}

type DeleteProductByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of product."`
}

// Integrate with Elasticsearch

type GetProductsWithElasticsearchRequest struct {
	Offset       int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit        int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy       string `query:"sort_by" default:"id:desc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. Format: \"field:asc/desc\" (default is asc if not declare after commas)"`
	Name         string `query:"name" example:"áo" doc:"Filter by name."`
	PriceGTE     string `query:"price_gte" pattern:"^[0-9]+$" example:"250000" doc:"Filter by price greater than or equal."`
	PriceLTE     string `query:"price_lte" pattern:"^[0-9]+$" example:"300000" doc:"Filter by price less than or equal."`
	CreatedAtGTE string `query:"created_at_gte" example:"2024-01-15T00:00:00" doc:"Filter by created_at greater than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
	CreatedAtLTE string `query:"created_at_lte" example:"2024-02-05T23:59:59" doc:"Filter by created_at less than or equal, with format is YYYY-MM-ddTHH:mm:ss."`
}
