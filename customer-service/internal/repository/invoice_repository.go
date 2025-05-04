package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type invoiceRepository struct {
}

type InvoiceRepository interface {
	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Invoice, error)
	GetById(ctx context.Context, id int64) (*model.Invoice, error)
	GetByUserId(ctx context.Context, userId int64, offset int, limit int, sortFields []utils.SortField) ([]model.Invoice, error)
	Create(ctx context.Context, newInvoice *model.Invoice) error
	UpdateById(ctx context.Context, id int64, updatedInvoice *model.Invoice) error
	DeleteById(ctx context.Context, id int64) error
}

func NewInvoiceRepository() InvoiceRepository {
	return &invoiceRepository{}
}

func (invoiceRepository *invoiceRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Invoice, error) {
	var invoices []model.Invoice
	query := infrastructure.DB.NewSelect().Model(&invoices).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (invoiceRepository *invoiceRepository) GetById(ctx context.Context, id int64) (*model.Invoice, error) {
	var invoice model.Invoice
	err := infrastructure.DB.NewSelect().Model(&invoice).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (invoiceRepository *invoiceRepository) GetByUserId(ctx context.Context, userId int64, offset int, limit int, sortFields []utils.SortField) ([]model.Invoice, error) {
	var invoices []model.Invoice
	query := infrastructure.DB.NewSelect().Model(&invoices).Where("user_id = ?", userId).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}

func (invoiceRepository *invoiceRepository) Create(ctx context.Context, newInvoice *model.Invoice) error {
	_, err := infrastructure.DB.NewInsert().Model(newInvoice).Exec(ctx)
	return err
}

func (invoiceRepository *invoiceRepository) UpdateById(ctx context.Context, id int64, updatedInvoice *model.Invoice) error {
	_, err := infrastructure.DB.NewUpdate().Model(updatedInvoice).Where("id = ?", id).Exec(ctx)
	return err
}

func (invoiceRepository *invoiceRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.DB.NewDelete().Model(&model.Invoice{}).Where("id = ?", id).Exec(ctx)
	return err
}
