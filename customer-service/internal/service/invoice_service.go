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
	invoiceRepository              repository.InvoiceRepository
	invoiceElasticsearchRepository repository.InvoiceElasticsearchRepository
}

type InvoiceService interface {
	GetInvoices(ctx context.Context, reqDTO *dto.GetInvoicesWithQueryParamRequest) ([]model.Invoice, error)
	GetInvoiceById(ctx context.Context, reqDTO *dto.GetInvoiceByIdRequest) (*model.Invoice, error)
	GetInvoicesByUserId(ctx context.Context, reqDTO *dto.GetInvoicesByUserIdWithQueryParamRequest) ([]model.Invoice, error)
	UpdateInvoiceById(ctx context.Context, reqDTO *dto.UpdateInvoiceRequest) error
	DeleteInvoiceById(ctx context.Context, reqDTO *dto.DeleteInvoiceRequest) error

	SyncAllInvoicesToElasticsearch(ctx context.Context) error

	GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) ([]model.Invoice, error)
	SumInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.AggregateInvoicesWithElasticsearchRequest) (*float64, error)
	SumAvgInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.AggregateInvoicesWithElasticsearchRequest) (*model.InvoiceReport, error)
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, invoiceElasticsearchRepository repository.InvoiceElasticsearchRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository:              invoiceRepository,
		invoiceElasticsearchRepository: invoiceElasticsearchRepository,
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

func (invoiceService *invoiceService) SyncAllInvoicesToElasticsearch(ctx context.Context) error {
	invoices, err := invoiceService.invoiceRepository.GetAll(ctx)

	if err != nil {
		return err
	}

	if err := invoiceService.invoiceElasticsearchRepository.SyncAll(ctx, invoices); err != nil {
		return err
	}

	return nil
}

func (invoiceService *invoiceService) GetInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.GetInvoicesWithElasticsearchRequest) ([]model.Invoice, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	invoices, err := invoiceService.invoiceElasticsearchRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields,
		reqDTO.CreatedAtGTE, reqDTO.CreatedAtLTE)
	if err != nil {
		return nil, err
	}

	return invoices, nil
}

func (invoiceService *invoiceService) SumInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.AggregateInvoicesWithElasticsearchRequest) (*float64, error) {
	sum, err := invoiceService.invoiceElasticsearchRepository.Sum(ctx, reqDTO.CreatedAtGTE, reqDTO.CreatedAtLTE)
	if err != nil {
		return nil, err
	}

	return sum, nil
}

func (invoiceService *invoiceService) SumAvgInvoicesWithElasticsearch(ctx context.Context, reqDTO *dto.AggregateInvoicesWithElasticsearchRequest) (*model.InvoiceReport, error) {
	invoiceReport, err := invoiceService.invoiceElasticsearchRepository.SumAvg(ctx, reqDTO.CreatedAtGTE, reqDTO.CreatedAtLTE)
	if err != nil {
		return nil, err
	}

	return invoiceReport, nil
}
