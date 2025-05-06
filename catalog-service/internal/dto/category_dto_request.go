package dto

type GetCategoriesRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:desc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. Format: \"field:asc/desc\" (default is asc if not declare after commas)"`
}

type GetCategoryByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of category."`
}

type GetCategoryByNameRequest struct {
	Name string `path:"name" required:"true" doc:"Name of category."`
}

type CreateCategoryRequest struct {
	Body struct {
		Name string `json:"name" required:"true" minLength:"1" doc:"Name of category."`
	}
}

type UpdateCategoryByIdRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of category will be updated."`
	Body struct {
		Name *string `json:"name,omitempty" minLength:"1" doc:"Name of category."`
	}
}

type DeleteCategoryByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of category."`
}
