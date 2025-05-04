package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type InvoiceDetailHandler struct {
	invoiceDetailService service.InvoiceDetailService
	invoiceSerivce       service.InvoiceService
	authMiddleware       *middleware.AuthMiddleware
}

func NewInvoiceDetailHandler(api huma.API, invoiceDetailService service.InvoiceDetailService, invoiceSerivce service.InvoiceService, authMiddleware *middleware.AuthMiddleware) *InvoiceDetailHandler {
	invoiceDetailHandler := &InvoiceDetailHandler{
		invoiceDetailService: invoiceDetailService,
		invoiceSerivce:       invoiceSerivce,
		authMiddleware:       authMiddleware,
	}

	// Get invoice details
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoice-details",
		Summary:     "/invoice-details",
		Description: "Get invoice details.",
		Tags:        []string{"Invoice Detail"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceDetailHandler.GetInvoiceDetails)

	// Get invoice detail by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoice-details/id/{id}",
		Summary:     "/invoice-details/id/{id}",
		Description: "Get invoice detail by id.",
		Tags:        []string{"Invoice Detail"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceDetailHandler.GetInvoiceDetailById)

	// Get invoice details by invoice id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/invoice-details/invoice-id/{invoice_id}",
		Summary:     "/invoice-details/invoice-id/{invoice_id}",
		Description: "Get invoice details by invoice id.",
		Tags:        []string{"Invoice Detail"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, invoiceDetailHandler.GetInvoiceDetailsByInvoiceId)

	// Get invoice details using account

	// Get invoice detail by id using account

	// Get invoice details by invoice id using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-invoice-details/invoice-id/{invoice_id}",
		Summary:     "/my-invoice-details/invoice-id/{invoice_id}",
		Description: "Get invoice details by invoice id using account.",
		Tags:        []string{"Invoice Detail"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, invoiceDetailHandler.GetInvoiceDetailsByInvoiceIdUsingAccount)

	return invoiceDetailHandler
}

func (invoiceDetailHandler *InvoiceDetailHandler) GetInvoiceDetails(ctx context.Context, reqDTO *dto.GetInvoiceDetailsWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.InvoiceDetailView], error) {
	invoiceDetails, err := invoiceDetailHandler.invoiceDetailService.GetInvoiceDetails(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoice details failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListInvoiceDetailView(invoiceDetails)
	res := &dto.PaginationBodyResponseList[dto.InvoiceDetailView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice details successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (invoiceDetailHandler *InvoiceDetailHandler) GetInvoiceDetailById(ctx context.Context, reqDTO *dto.GetInvoiceDetailByIdRequest) (*dto.BodyResponse[dto.InvoiceDetailView], error) {
	foundInvoiceDetail, err := invoiceDetailHandler.invoiceDetailService.GetInvoiceDetailById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice detail by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToInvoiceDetailView(foundInvoiceDetail)
	res := &dto.BodyResponse[dto.InvoiceDetailView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice detail by id successful"
	res.Body.Data = *data
	return res, nil
}

func (invoiceDetailHandler *InvoiceDetailHandler) GetInvoiceDetailsByInvoiceId(ctx context.Context, reqDTO *dto.GetInvoiceDetailsByInvoiceIdWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.InvoiceDetailView], error) {
	invoiceDetails, err := invoiceDetailHandler.invoiceDetailService.GetInvoiceDetailsByInvoiceId(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoice details by invoice id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListInvoiceDetailView(invoiceDetails)
	res := &dto.PaginationBodyResponseList[dto.InvoiceDetailView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice details by invoice id successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (invoiceDetailHandler *InvoiceDetailHandler) GetInvoiceDetailsByInvoiceIdUsingAccount(ctx context.Context, reqDTO *dto.GetInvoiceDetailsByInvoiceIdUsingAccountWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.InvoiceDetailView], error) {
	userId := ctx.Value("user_id").(int64)

	foundInvoice, err := invoiceDetailHandler.invoiceSerivce.GetInvoiceById(ctx, &dto.GetInvoiceByIdRequest{Id: reqDTO.InvoiceId})
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice details by invoice id using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	} else if foundInvoice.UserId != userId {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get invoice details by invoice id using account failed"
		res.Details = []string{"id of invoice is not valid"}
		return nil, res
	}

	convertReqDTO := &dto.GetInvoiceDetailsByInvoiceIdWithQueryParamRequest{
		InvoiceId: reqDTO.InvoiceId,
		Offset:    reqDTO.Offset,
		Limit:     reqDTO.Limit,
		SortBy:    reqDTO.SortBy,
	}

	invoiceDetails, err := invoiceDetailHandler.invoiceDetailService.GetInvoiceDetailsByInvoiceId(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get invoice details by invoice id using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListInvoiceDetailView(invoiceDetails)
	res := &dto.PaginationBodyResponseList[dto.InvoiceDetailView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get invoice details by invoice id using account successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}
