package service

import "fmt"

func IssueInvoice(invoice Invoice) (Invoice, error) {
	if invoice.Status != InvoiceStatusDraft {
		return Invoice{}, fmt.Errorf("invoice status %q cannot transition to issued", invoice.Status)
	}

	invoice.Status = InvoiceStatusIssued
	invoice.Items = append([]LineItem(nil), invoice.Items...)
	return invoice, nil
}
