package dto

import (
	"thanhldt060802/internal/model"

	"github.com/shopspring/decimal"
)

type InvoiceDetailView struct {
	Id                 int64           `json:"id"`
	InvoiceId          int64           `json:"invoice_id"`
	ProductId          int64           `json:"product_id"`
	Price              decimal.Decimal `json:"price"`
	DiscountPercentage int32           `json:"discount_percentage"`
	Quantity           int32           `json:"quantity"`
	TotalPrice         decimal.Decimal `json:"total_price"`
}

func ToInvoiceDetailView(invoiceDetail *model.InvoiceDetail) *InvoiceDetailView {
	return &InvoiceDetailView{
		Id:                 invoiceDetail.Id,
		InvoiceId:          invoiceDetail.InvoiceId,
		ProductId:          invoiceDetail.ProductId,
		Price:              invoiceDetail.Price,
		DiscountPercentage: invoiceDetail.DiscountPercentage,
		Quantity:           invoiceDetail.Quantity,
		TotalPrice:         invoiceDetail.TotalPrice,
	}
}

func ToListInvoiceDetailView(invoiceDetails []model.InvoiceDetail) []InvoiceDetailView {
	invoiceDetailViews := make([]InvoiceDetailView, len(invoiceDetails))
	for i, invoiceDetail := range invoiceDetails {
		invoiceDetailViews[i] = *ToInvoiceDetailView(&invoiceDetail)
	}
	return invoiceDetailViews
}
