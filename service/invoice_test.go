package service

import "testing"

func TestNewDraftInvoiceCalculatesTotals(t *testing.T) {
	items := []LineItem{
		{Quantity: 2, UnitPriceMinor: 1500, TaxRateBPS: 500},
	}
	invoice, err := NewDraftInvoice("customer-1", "USD", items)
	if err != nil {
		t.Fatal(err)
	}

	if invoice.CustomerID != "customer-1" {
		t.Fatalf("CustomerID = %q, want %q", invoice.CustomerID, "customer-1")
	}
	if invoice.Currency != "USD" {
		t.Fatalf("Currency = %q, want %q", invoice.Currency, "USD")
	}
	if invoice.Status != InvoiceStatusDraft {
		t.Fatalf("Status = %q, want %q", invoice.Status, InvoiceStatusDraft)
	}
	if invoice.Totals.TotalMinor != 3150 {
		t.Fatalf("TotalMinor = %d, want %d", invoice.Totals.TotalMinor, 3150)
	}

	items[0].UnitPriceMinor = 9999
	if invoice.Items[0].UnitPriceMinor != 1500 {
		t.Fatal("invoice items changed after input mutation")
	}
}

func TestNewDraftInvoiceRejectsMissingIdentity(t *testing.T) {
	if _, err := NewDraftInvoice("", "USD", []LineItem{{Quantity: 1, UnitPriceMinor: 100}}); err == nil {
		t.Fatal("NewDraftInvoice() succeeded without customer ID")
	}
	if _, err := NewDraftInvoice("   ", "USD", []LineItem{{Quantity: 1, UnitPriceMinor: 100}}); err == nil {
		t.Fatal("NewDraftInvoice() succeeded with whitespace customer ID")
	}
	if _, err := NewDraftInvoice("customer-1", "USD", nil); err == nil {
		t.Fatal("NewDraftInvoice() succeeded without line items")
	}
}

func TestNewDraftInvoiceRejectsMalformedCurrency(t *testing.T) {
	for _, currency := range []string{"US", "usd", "123", "!!!", "ÄBC"} {
		if _, err := NewDraftInvoice("customer-1", currency, []LineItem{{Quantity: 1, UnitPriceMinor: 100}}); err == nil {
			t.Errorf("NewDraftInvoice() accepted malformed currency %q", currency)
		}
	}
}
