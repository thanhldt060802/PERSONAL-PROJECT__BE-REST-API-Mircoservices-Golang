package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type invoiceDetailRepository struct {
}

type InvoiceDetailRepository interface {
	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.InvoiceDetail, error)
	GetById(ctx context.Context, id int64) (*model.InvoiceDetail, error)
	GetByInvoiceId(ctx context.Context, invoiceId int64, offset int, limit int, sortFields []utils.SortField) ([]model.InvoiceDetail, error)
	Create(ctx context.Context, newInvoiceDetail *model.InvoiceDetail) error
	UpdateById(ctx context.Context, id int64, updatedInvoiceDetail *model.InvoiceDetail) error
	DeleteById(ctx context.Context, id int64) error
}

func NewInvoiceDetailRepository() InvoiceDetailRepository {
	return &invoiceDetailRepository{}
}

func (invoiceDetailRepository *invoiceDetailRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.InvoiceDetail, error) {
	var invoiceDetails []model.InvoiceDetail
	query := infrastructure.DB.NewSelect().Model(&invoiceDetails).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return invoiceDetails, nil
}

func (invoiceDetailRepository *invoiceDetailRepository) GetById(ctx context.Context, id int64) (*model.InvoiceDetail, error) {
	var invoiceDetail model.InvoiceDetail
	err := infrastructure.DB.NewSelect().Model(&invoiceDetail).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &invoiceDetail, nil
}

func (invoiceDetailRepository *invoiceDetailRepository) GetByInvoiceId(ctx context.Context, invoiceId int64, offset int, limit int, sortFields []utils.SortField) ([]model.InvoiceDetail, error) {
	var invoiceDetails []model.InvoiceDetail
	query := infrastructure.DB.NewSelect().Model(&invoiceDetails).Where("invoice_id = ?", invoiceId).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return invoiceDetails, nil
}

func (invoiceDetailRepository *invoiceDetailRepository) Create(ctx context.Context, newInvoiceDetail *model.InvoiceDetail) error {
	_, err := infrastructure.DB.NewInsert().Model(newInvoiceDetail).Exec(ctx)
	return err
}

func (invoiceDetailRepository *invoiceDetailRepository) UpdateById(ctx context.Context, id int64, updatedInvoiceDetail *model.InvoiceDetail) error {
	_, err := infrastructure.DB.NewUpdate().Model(updatedInvoiceDetail).Where("id = ?", id).Exec(ctx)
	return err
}

func (invoiceDetailRepository *invoiceDetailRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.DB.NewDelete().Model(&model.InvoiceDetail{}).Where("id = ?", id).Exec(ctx)
	return err
}
