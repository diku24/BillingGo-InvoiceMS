package service

import "testing"

func TestCalculateTotals(t *testing.T) {
	items := []LineItem{
		{Quantity: 2, UnitPriceMinor: 1000, TaxRateBPS: 1000},
		{Quantity: 1, UnitPriceMinor: 500, TaxRateBPS: 0, DiscountMinor: 100},
	}

	got, err := CalculateTotals(items)
	if err != nil {
		t.Fatal(err)
	}

	want := Totals{
		SubtotalMinor: 2500,
		DiscountMinor: 100,
		TaxMinor:      200,
		TotalMinor:    2600,
	}
	if got != want {
		t.Fatalf("CalculateTotals() = %+v, want %+v", got, want)
	}
}

func TestCalculateTotalsRejectsInvalidLineItems(t *testing.T) {
	tests := []LineItem{
		{Quantity: 0, UnitPriceMinor: 100},
		{Quantity: 1, UnitPriceMinor: -1},
		{Quantity: 1, UnitPriceMinor: 100, TaxRateBPS: -1},
		{Quantity: 1, UnitPriceMinor: 100, DiscountMinor: -1},
		{Quantity: 1, UnitPriceMinor: 100, DiscountMinor: 101},
	}

	for _, item := range tests {
		if _, err := CalculateTotals([]LineItem{item}); err == nil {
			t.Errorf("CalculateTotals(%+v) succeeded, want validation error", item)
		}
	}
}

func TestCalculateTotalsRejectsArithmeticOverflow(t *testing.T) {
	const maxInt64 = int64(1<<63 - 1)

	tests := []LineItem{
		{Quantity: 2, UnitPriceMinor: maxInt64},
		{Quantity: 1, UnitPriceMinor: maxInt64, TaxRateBPS: 1},
	}

	for _, item := range tests {
		if _, err := CalculateTotals([]LineItem{item}); err == nil {
			t.Errorf("CalculateTotals(%+v) succeeded, want overflow error", item)
		}
	}

	items := []LineItem{
		{Quantity: 1, UnitPriceMinor: maxInt64 / 2},
		{Quantity: 1, UnitPriceMinor: maxInt64 / 2},
		{Quantity: 1, UnitPriceMinor: 2},
	}
	if _, err := CalculateTotals(items); err == nil {
		t.Fatal("CalculateTotals() succeeded with aggregate overflow")
	}
}
