package main

import (
	"net/http"
	"thanhldt060802/config"
	"thanhldt060802/database"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/handler"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/repository"
	"thanhldt060802/internal/service"
	"thanhldt060802/redis"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

// Huma Docs UI template by Scalar
var humaDocsEmbedded = `<!doctype html>
<html>
  <head>
    <title>API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`

func main() {

	config.InitConfig()

	db := database.ConnectDB()
	defer db.Close()

	r := gin.Default()

	humaCfg := huma.DefaultConfig("Catalog Service", "v1.0.0")
	humaCfg.DocsPath = ""
	humaCfg.JSONSchemaDialect = ""
	humaCfg.CreateHooks = nil
	humaCfg.Components = &huma.Components{
		SecuritySchemes: map[string]*huma.SecurityScheme{
			"BearerAuth": {
				Type:         "http",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
	}

	api := humagin.New(r, humaCfg)

	r.GET("/docs", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html", []byte(humaDocsEmbedded))
	})

	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		details := make([]string, len(errs))
		for i, err := range errs {
			details[i] = err.Error()
		}
		res := &dto.ErrorResponse{}
		res.Status = status
		res.Message = msg
		res.Error_ = "Some thing went wrong"
		res.Details = details
		return res
	}

	// Initialize Redis client
	redisClient := redis.NewClient()

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(api, redisClient)

	// Initialize repositories
	userRepository := repository.NewProductRepository(db)

	// Initialize services
	userService := service.NewProductService(userRepository)

	// Initialize handlers
	handler.NewProductHandler(api, userService, authMiddleware)

	r.Run(":" + config.AppConfig.AppPort)

}
