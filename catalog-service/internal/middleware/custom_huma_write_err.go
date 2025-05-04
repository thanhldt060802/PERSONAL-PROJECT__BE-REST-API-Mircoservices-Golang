package middleware

import (
	"encoding/json"
	"thanhldt060802/internal/dto"

	"github.com/danielgtaylor/huma/v2"
)

func CustomerHumaWriteErr(ctx huma.Context, status int, code string, message string, details []string) {
	res := &dto.ErrorResponse{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}

	ctx.SetHeader("Content-Type", "application/json")
	ctx.SetStatus(status)
	json.NewEncoder(ctx.BodyWriter()).Encode(res)
}
