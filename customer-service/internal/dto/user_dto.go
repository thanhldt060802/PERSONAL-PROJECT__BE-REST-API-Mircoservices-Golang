package dto

import (
	"strings"
	"thanhldt060802/internal/model"
	"time"
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
type UserDTO struct {
	Id        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Address   string    `json:"address"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserDTO(user *model.User) *UserDTO {
	return &UserDTO{
		Id:        user.Id,
		FullName:  user.FullName,
		Email:     user.Email,
		Username:  user.Username,
		Address:   user.Address,
		RoleName:  user.RoleName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserDTOs(users []model.User) []UserDTO {
	userDTOs := make([]UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = *ToUserDTO(&user)
	}
	return userDTOs
}

// ################################################################################

// Request
// ################################################################################
type GetUsersRequestQueryParam struct {
	Offset int    `query:"offset" default:"0" minimum:"0" doc:"Skip item by offset."`
	Limit  int    `query:"limit" default:"5" minimum:"1" maximum:"10" doc:"Limit item from offset."`
	SortBy string `query:"sort_by" default:"full_name:asc" example:"full_name,created_at:desc" doc:"Sort by one or more fields separated by commas. For example: sort_by=name,created_at:desc will sort by name in ascending order, then by created_at in descending order."`
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
	Id   int64 `path:"id" required:"true"`
	Body struct {
		FullName *string `json:"fullname,omitempty" minLength:"1" doc:"Full name of user account."`
		Email    *string `json:"email,omitempty" minLength:"1" format:"email" doc:"Email of user account."`
		Password *string `json:"password,omitempty" minLength:"1" doc:"Password of user account."`
		Address  *string `json:"address,omitempty" minLength:"1" doc:"Address of user account."`
		RoleName *string `json:"role_name,omitempty" enum:"ADMIN,CUSTOMER" doc:"Role name of user account."`
	}
}

type DeleteUserRequest struct {
	Id int64 `path:"id" required:"true"`
}

type LoginRequest struct {
	Body struct {
		Username string `json:"username" required:"true" minLength:"1" example:"user1" doc:"Account username."`
		Password string `json:"password" required:"true" minLength:"1" example:"123" doc:"Account password."`
	}
}

type RegisterRequest struct {
	Body struct {
		FullName string `json:"full_name" required:"true" minLength:"1" doc:"Full name of user account."`
		Email    string `json:"email" required:"true" minLength:"1" doc:"Email of user account."`
		Username string `json:"username" required:"true" minLength:"1" doc:"Username of user account."`
		Password string `json:"password" required:"true" minLength:"1" doc:"Password of user account."`
		Address  string `json:"address" required:"true" minLength:"1" doc:"Address of user account."`
	}
}

type UpdateAccountInfoRequest struct {
	Body struct {
		FullName *string `json:"fullname,omitempty" minLength:"1" doc:"Full name of user account."`
		Email    *string `json:"email,omitempty" minLength:"1" format:"email" doc:"Email of user account."`
		Password *string `json:"password,omitempty" minLength:"1" doc:"Password of user account."`
		Address  *string `json:"address,omitempty" minLength:"1" doc:"Address of user account."`
	}
}

// ################################################################################
