package dto

import (
	"thanhldt060802/internal/model"
	"time"

	"github.com/shopspring/decimal"
)

type InvoiceView struct {
	Id          int64           `json:"id"`
	UserId      int64           `json:"user_id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Stautus     string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func ToInvoiceView(invoice *model.Invoice) *InvoiceView {
	return &InvoiceView{
		Id:          invoice.Id,
		UserId:      invoice.UserId,
		TotalAmount: invoice.TotalAmount,
		Stautus:     invoice.Status,
		CreatedAt:   invoice.CreatedAt,
		UpdatedAt:   invoice.UpdatedAt,
	}
}

func ToListInvoiceView(invoices []model.Invoice) []InvoiceView {
	invoiceViews := make([]InvoiceView, len(invoices))
	for i, invoice := range invoices {
		invoiceViews[i] = *ToInvoiceView(&invoice)
	}
	return invoiceViews
}
