package service

import (
	"fmt"
	"strings"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusIssued    InvoiceStatus = "issued"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	CustomerID string
	Currency   string
	Items      []LineItem
	Totals     Totals
	Status     InvoiceStatus
}

func isCurrencyCode(currency string) bool {
	if len(currency) != 3 {
		return false
	}
	for _, character := range currency {
		if character < 'A' || character > 'Z' {
			return false
		}
	}
	return true
}

func NewDraftInvoice(customerID, currency string, items []LineItem) (Invoice, error) {
	if strings.TrimSpace(customerID) == "" {
		return Invoice{}, fmt.Errorf("customer ID is required")
	}
	if !isCurrencyCode(currency) {
		return Invoice{}, fmt.Errorf("currency must be a three-letter uppercase code")
	}
	if len(items) == 0 {
		return Invoice{}, fmt.Errorf("at least one line item is required")
	}

	totals, err := CalculateTotals(items)
	if err != nil {
		return Invoice{}, fmt.Errorf("calculate invoice totals: %w", err)
	}

	invoiceItems := append([]LineItem(nil), items...)
	return Invoice{
		CustomerID: customerID,
		Currency:   currency,
		Items:      invoiceItems,
		Totals:     totals,
		Status:     InvoiceStatusDraft,
	}, nil
}
