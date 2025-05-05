package dto

import "github.com/shopspring/decimal"

// Only user request
// ################################################################################

type GetUsersWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetUserByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of user will be gotten."`
}

type GetUserByUsernameRequest struct {
	Username string `path:"username" required:"true" doc:"Username of user will be gotten."`
}

type GetUserByEmailRequest struct {
	Email string `path:"email" required:"true" doc:"Email of user will be gotten."`
}

type CreateUserRequest struct {
	Body struct {
		FullName string `json:"full_name" required:"true" minLength:"1" doc:"Full name of user acount."`
		Email    string `json:"email" required:"true" format:"email" doc:"Email of user acount."`
		Username string `json:"username" required:"true" minLength:"1" doc:"Username of user acount."`
		Password string `json:"password" required:"true" minLength:"1" doc:"Password of user acount."`
		Address  string `json:"address" required:"true" minLength:"1" doc:"Address of user acount."`
		RoleName string `json:"role_name" required:"true" enum:"ADMIN,CUSTOMER" doc:"Role name of user account."`
	}
}

type UpdateUserRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of user will be updated."`
	Body struct {
		FullName *string `json:"fullname,omitempty" minLength:"1" doc:"Full name of user account."`
		Email    *string `json:"email,omitempty" minLength:"1" format:"email" doc:"Email of user account."`
		Password *string `json:"password,omitempty" minLength:"1" doc:"Password of user account."`
		Address  *string `json:"address,omitempty" minLength:"1" doc:"Address of user account."`
		RoleName *string `json:"role_name,omitempty" enum:"ADMIN,CUSTOMER" doc:"Role name of user account."`
	}
}

type DeleteUserRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of user will be deleted."`
}

type LoginUserRequest struct {
	Body struct {
		Username string `json:"username" required:"true" minLength:"1" example:"user1" doc:"Account username."`
		Password string `json:"password" required:"true" minLength:"1" example:"123" doc:"Account password."`
	}
}

type LogoutUserRequest struct {
	Body struct {
		Token string `json:"token" required:"true" minLength:"1" example:"XXX" doc:"Token of account will be logged out."`
	}
}

type RegisterRequest struct {
	Body struct {
		FullName string `json:"full_name" required:"true" minLength:"1" doc:"Full name of user acount."`
		Email    string `json:"email" required:"true" format:"email" doc:"Email of user acount."`
		Username string `json:"username" required:"true" minLength:"1" doc:"Username of user acount."`
		Password string `json:"password" required:"true" minLength:"1" doc:"Password of user acount."`
		Address  string `json:"address" required:"true" minLength:"1" doc:"Address of user acount."`
	}
}

type UpdateUserUsingAccountRequest struct {
	Body struct {
		FullName *string `json:"fullname,omitempty" minLength:"1" doc:"Full name of user account."`
		Email    *string `json:"email,omitempty" minLength:"1" format:"email" doc:"Email of user account."`
		Password *string `json:"password,omitempty" minLength:"1" doc:"Password of user account."`
		Address  *string `json:"address,omitempty" minLength:"1" doc:"Address of user account."`
	}
}

// ################################################################################

// Only cart request
// ################################################################################

type GetCartsWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetCartByUserIdRequest struct {
	UserId int64 `path:"user_id" required:"true" doc:"User id of cart will be gotten."`
}

// ################################################################################

// Only cart item request
// ################################################################################

type GetCartItemsWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetCartItemByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of cart item will be gotten."`
}

type GetCartItemsByCartIdWithQueryParamRequest struct {
	CartId int64  `path:"cart_id" required:"true" doc:"Cart id of cart items will be filtered."`
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
}

type CreateCartItemRequest struct {
	CartId int64 `path:"id" required:"true" doc:"Id of cart will be added cart item."`
	Body   struct {
		ProductId int64 `json:"product_id" required:"true" minimum:"1" doc:"Product id of cart item."`
	}
}

type UpdateCartItemRequest struct {
	CartId int64 `path:"cart_id" required:"true" doc:"Id of cart will be updated cart item."`
	Id     int64 `path:"id" required:"true" doc:"Id of cart item will be updated."`
	Body   struct {
		Quantity *int32 `json:"quantity,omitempty" minimum:"1" doc:"Quantity of cart item."`
	}
}

type DeleteCartItemRequest struct {
	CartId int64 `path:"cart_id" required:"true" doc:"Id of cart will be deleted cart item."`
	Id     int64 `path:"id" required:"true" doc:"Id of cart item will be deleted."`
}

type GetCartItemsUsingAccountWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
}

type CreateCartItemUsingAccountRequest struct {
	Body struct {
		ProductId int64 `json:"product_id" required:"true" minimum:"1" doc:"Product id of cart item."`
	}
}

type UpdateCartItemUsingAccountRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of cart item will be updated."`
	Body struct {
		Quantity *int32 `json:"quantity,omitempty" minimum:"1" doc:"Quantity of cart item."`
	}
}

type DeleteCartItemUsingAccountRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of cart item will be deleted."`
}

// ################################################################################

// Only invoice request
// ################################################################################

type GetInvoicesWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetInvoiceByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of invoice item will be gotten."`
}

type GetInvoicesByUserIdWithQueryParamRequest struct {
	UserId int64  `path:"user_id" required:"true" doc:"User id of invoices will be filtered."`
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
}

type CreateInvoiceRequest struct {
	Body struct {
		UserId      int64           `json:"user_id" required:"true" minimum:"1" doc:"User id of invoice."`
		TotalAmount decimal.Decimal `json:"total_amount" required:"true" minimum:"0" doc:"Total amount of invoice."`
		Status      string          `json:"status" required:"true" minimum:"1" doc:"Status of invoice."`
	}
}

type UpdateInvoiceRequest struct {
	Id   int64 `path:"id" required:"true" doc:"Id of invoice will be updated."`
	Body struct {
		Status *string `json:"status,omitempty" minimum:"1" doc:"Status of invoice."`
	}
}

type DeleteInvoiceRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of invoice will be deleted."`
}

type GetInvoicesUsingAccountQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetInvoiceByIdUsingAccountRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of invoice item will be gotten."`
}

type DeleteInvoiceUsingAccountRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of invoice will be deleted."`
}

// ################################################################################

// Only invoice detail request
// ################################################################################

type GetInvoiceDetailsWithQueryParamRequest struct {
	Offset int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"id:asc" example:"created_at:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=created_at:desc,id will sort by created_at in descending order, then by id in ascending order."`
}

type GetInvoiceDetailByIdRequest struct {
	Id int64 `path:"id" required:"true" doc:"Id of invoice detail will be gotten."`
}

type GetInvoiceDetailsByInvoiceIdWithQueryParamRequest struct {
	InvoiceId int64  `path:"invoice_id" required:"true" doc:"Invoice id of invoice details will be filtered."`
	Offset    int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit     int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy    string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
}

type GetInvoiceDetailsByInvoiceIdUsingAccountWithQueryParamRequest struct {
	InvoiceId int64  `path:"invoice_id" required:"true" doc:"Invoice id of invoice details will be filtered."`
	Offset    int    `query:"offset" default:"0" minimum:"0" example:"0" doc:"Skip item by offset."`
	Limit     int    `query:"limit" default:"5" minimum:"1" maximum:"10" example:"10" doc:"Limit item from offset."`
	SortBy    string `query:"sort_by" default:"id:asc" example:"quantity:desc,id" doc:"Sort by one or more fields separated by commas. For example: sort_by=quantity:desc,id will sort by quantity in descending order, then by id in ascending order."`
}

// ################################################################################
