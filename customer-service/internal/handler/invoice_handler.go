package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type InvoiceHandler struct {
	invoiceService service.InvoiceService
	authMiddleware *middleware.AuthMiddleware
}

func NewInvoiceHandler(api huma.API, invoiceService service.InvoiceService, authMiddleware *middleware.AuthMiddleware) *InvoiceHandler {
	invoiceHandler := &InvoiceHandler{
		invoiceService: invoiceService,
		authMiddleware: authMiddleware,
	}

	// Get invoices
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices",
		Summary:     "/invoices",
		Description: "Get invoices.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoices)

	// Get invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/id/{id}",
		Summary:     "/invoices/id/{id}",
		Description: "Get invoice by id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoiceById)

	// Get invoices by user id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/user-id/{user_id}",
		Summary:     "/invoices/user-id/{user_id}",
		Description: "Get invoices by user id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceHandler.GetInvoicesByUserId)

	// Update invoice by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/invoices/id/{id}",
		Summary:     "/invoices/id/{id}",
		Description: "Update invoice by id.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceHandler.UpdateInvoiceById)

	// Get invoices using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-invoices",
		Summary:     "/my-invoices",
		Description: "Get invoices using account.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, invoiceHandler.GetInvoicesUsingAccount)

	// Get invoice by id using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-invoices/id/{id}",
		Summary:     "/my-invoices/id/{id}",
		Description: "Get invoice by id using account.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, invoiceHandler.GetInvoiceByIdUsingAccount)

	// Delete invoice by id using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/my-invoices/id/{id}",
		Summary:     "/my-invoices/id/{id}",
		Description: "Delete invoice by id using account.",
		Tags:        []string{"Invoice"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceHandler.DeleteInvoiceByIdUsingAccount)

	// Sync all invoices to Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/sync-to-elasticsearch",
		Summary:     "/invoices/sync-to-elasticsearch",
		Description: "Sync all invoices to Elasticsearch.",
		Tags:        []string{"Invoice"},
		// Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceHandler.SyncAllInvoicesToElasticsearch)

	// Get invoices with Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/elasticsearch",
		Summary:     "/invoices/elasticsearch",
		Description: "Get invoices with Elasticsearch.",
		Tags:        []string{"Invoice"},
	}, invoiceHandler.GetInvoicesWithElasticsearch)

	// Sum invoices with Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/elasticsearch/sum",
		Summary:     "/invoices/elasticsearch/sum",
		Description: "Sum invoices with Elasticsearch.",
		Tags:        []string{"Invoice"},
	}, invoiceHandler.SumInvoicesWithElasticsearch)

	// Sum avg invoices with Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoices/elasticsearch/sumavg",
		Summary:     "/invoices/elasticsearch/sumavg",
		Description: "Sum avg invoices with Elasticsearch.",
		Tags:        []string{"Invoice"},
	}, invoiceHandler.SumAvgInvoicesWithElasticsearch)

	return invoiceHandler
}

func (invoiceHandler *InvoiceHandler) GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoices(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListInvoiceView(invoices)
	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*dto.BodyResponse[dto.InvoiceView], error) {
	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToInvoiceView(foundInvoice)
	res := &dto.BodyResponse[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice by id successful"
	res.Body.Data = *data
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoicesByUserId(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListInvoiceView(invoices)
	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceRequest) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.UpdateInvoiceById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update invoice failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update invoice successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceRequest) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.DeleteInvoiceById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete invoice failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete invoice successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoicesUsingAccount(ctx context.Context, reqDTO *dto.GetInvoicesUsingAccountQueryParamRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	userId := ctx.Value("user_id").(int64)

	convertReqDTO := &dto.GetInvoicesByUserIdWithQueryParamRequest{
		UserId: userId,
		Offset: reqDTO.Offset,
		Limit:  reqDTO.Limit,
		SortBy: reqDTO.SortBy,
	}

	invoices, err := invoiceHandler.invoiceService.GetInvoicesByUserId(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListInvoiceView(invoices)
	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoiceByIdUsingAccount(ctx context.Context, reqDTO *dto.GetInvoiceByIdUsingAccountRequest) (*dto.BodyResponse[dto.InvoiceView], error) {
	userId := ctx.Value("user_id").(int64)

	convertReqDTO := &dto.GetInvoiceByIdRequest{Id: reqDTO.Id}

	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice by id using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	} else if foundInvoice.UserId != userId {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice by id using account failed"
		res.Details = []string{"id of invoice is not valid"}
		return nil, res
	}

	data := dto.ToInvoiceView(foundInvoice)
	res := &dto.BodyResponse[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice by id using account successful"
	res.Body.Data = *data
	return res, nil
}

func (invoiceHandler *InvoiceHandler) DeleteInvoiceByIdUsingAccount(ctx context.Context, reqDTO *dto.DeleteInvoiceUsingAccountRequest) (*dto.SuccessResponse, error) {
	userId := ctx.Value("user_id").(int64)

	foundInvoice, err := invoiceHandler.invoiceService.GetInvoiceById(ctx, &dto.GetInvoiceByIdRequest{Id: reqDTO.Id})
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete invoice using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	} else if foundInvoice.UserId != userId {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete invoice using account failed"
		res.Details = []string{"id of invoice is not valid"}
		return nil, res
	}

	convertReqDTO := &dto.DeleteInvoiceRequest{Id: reqDTO.Id}

	if err := invoiceHandler.invoiceService.DeleteInvoiceById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete invoice using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete invoice using account successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) SyncAllInvoicesToElasticsearch(ctx context.Context, reqDTO *struct{}) (*dto.SuccessResponse, error) {
	if err := invoiceHandler.invoiceService.SyncAllInvoicesToElasticsearch(ctx); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sync all invoices to Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Sync all invoices to Elasticsearch successful"
	return res, nil
}

func (invoiceHandler *InvoiceHandler) GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) (*dto.PaginationBodyResponseList[dto.InvoiceView], error) {
	invoices, err := invoiceHandler.invoiceService.GetInvoicesWithElasticsearch(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoices with Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListInvoiceView(invoices)
	res := &dto.PaginationBodyResponseList[dto.InvoiceView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoices with Elasticsearch successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (invoiceHandler *InvoiceHandler) SumInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.AggregateInvoicesWithElasticsearchRequest) (*dto.BodyResponse[float64], error) {
	sum, err := invoiceHandler.invoiceService.SumInvoicesWithElasticsearch(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sum invoices with Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[float64]{}
	res.Body.Code = "OK"
	res.Body.Message = "Sum invoices with Elasticsearch successful"
	res.Body.Data = *sum
	return res, nil
}

func (invoiceHandler *InvoiceHandler) SumAvgInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.AggregateInvoicesWithElasticsearchRequest) (*dto.BodyResponse[model.InvoiceReport], error) {
	invoiceReport, err := invoiceHandler.invoiceService.SumAvgInvoicesWithElasticsearch(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sum avg invoices with Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[model.InvoiceReport]{}
	res.Body.Code = "OK"
	res.Body.Message = "Sum avg invoices with Elasticsearch successful"
	res.Body.Data = *invoiceReport
	return res, nil
}
