package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"
)

type invoiceService struct {
	invoiceRepository repository.InvoiceRepository
}

type InvoiceService interface {
	GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesWithQueryParamRequest) ([]model.Invoice, error)
	GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*model.Invoice, error)
	GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdWithQueryParamRequest) ([]model.Invoice, error)
	UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceRequest) error
	DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceRequest) error
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository: invoiceRepository,
	}
}

func (invoiceService *invoiceService) GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesWithQueryParamRequest) ([]model.Invoice, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	invoices, err := invoiceService.invoiceRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (invoiceService *invoiceService) GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*model.Invoice, error) {
	foundInvoice, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, err
	}

	return foundInvoice, nil
}

func (invoiceService *invoiceService) GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdWithQueryParamRequest) ([]model.Invoice, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	invoices, err := invoiceService.invoiceRepository.GetByUserId(ctx, reqDTO.UserId, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (invoiceService *invoiceService) UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceRequest) error {
	foundInvoice, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of invoice is not valid")
	}

	if reqDTO.Body.Status != nil {
		foundInvoice.Status = *reqDTO.Body.Status
	}
	foundInvoice.UpdatedAt = time.Now().UTC()

	if err := invoiceService.invoiceRepository.UpdateById(ctx, reqDTO.Id, foundInvoice); err != nil {
		return err
	}

	return nil
}

func (invoiceService *invoiceService) DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceRequest) error {
	if _, err := invoiceService.invoiceRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of invoice is not valid")
	}

	if err := invoiceService.invoiceRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return err
	}

	return nil
}
