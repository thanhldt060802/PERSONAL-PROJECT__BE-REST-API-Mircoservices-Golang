package service

import (
	"context"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
)

type invoiceDetailService struct {
	invoiceDetailRepository repository.InvoiceDetailRepository
}

type InvoiceDetailService interface {
	GetInvoiceDetails(ctx context.Context, reqDTO *dto.GetInvoiceDetailsWithQueryParamRequest) ([]model.InvoiceDetail, error)
	GetInvoiceDetailById(ctx context.Context, reqDTO *dto.GetInvoiceDetailByIdRequest) (*model.InvoiceDetail, error)
	GetInvoiceDetailsByInvoiceId(ctx context.Context, reqDTO *dto.GetInvoiceDetailsByInvoiceIdWithQueryParamRequest) ([]model.InvoiceDetail, error)
}

func NewInvoiceDetailService(invoiceDetailRepository repository.InvoiceDetailRepository) InvoiceDetailService {
	return &invoiceDetailService{
		invoiceDetailRepository: invoiceDetailRepository,
	}
}

func (invoiceDetailService *invoiceDetailService) GetInvoiceDetails(ctx context.Context, reqDTO *dto.GetInvoiceDetailsWithQueryParamRequest) ([]model.InvoiceDetail, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	invoiceDetails, err := invoiceDetailService.invoiceDetailRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return invoiceDetails, nil
}

func (invoiceDetailService *invoiceDetailService) GetInvoiceDetailById(ctx context.Context, reqDTO *dto.GetInvoiceDetailByIdRequest) (*model.InvoiceDetail, error) {
	foundInvoiceDetail, err := invoiceDetailService.invoiceDetailRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, err
	}

	return foundInvoiceDetail, nil
}

func (invoiceDetailService *invoiceDetailService) GetInvoiceDetailsByInvoiceId(ctx context.Context, reqDTO *dto.GetInvoiceDetailsByInvoiceIdWithQueryParamRequest) ([]model.InvoiceDetail, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	invoiceDetails, err := invoiceDetailService.invoiceDetailRepository.GetByInvoiceId(ctx, reqDTO.InvoiceId, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return invoiceDetails, nil
}
