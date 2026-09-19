package service

import "testing"

func TestIssueInvoiceTransitionsDraft(t *testing.T) {
	draft, err := NewDraftInvoice("customer-1", "USD", []LineItem{{Quantity: 1, UnitPriceMinor: 100}})
	if err != nil {
		t.Fatal(err)
	}

	issued, err := IssueInvoice(draft)
	if err != nil {
		t.Fatal(err)
	}
	if issued.Status != InvoiceStatusIssued {
		t.Fatalf("Status = %q, want %q", issued.Status, InvoiceStatusIssued)
	}
	if draft.Status != InvoiceStatusDraft {
		t.Fatal("IssueInvoice mutated the draft invoice")
	}
}

func TestIssueInvoiceRejectsNonDraft(t *testing.T) {
	invoice := Invoice{Status: InvoiceStatusIssued}
	if _, err := IssueInvoice(invoice); err == nil {
		t.Fatal("IssueInvoice() succeeded for an already-issued invoice")
	}
}
