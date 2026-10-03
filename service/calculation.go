package service

import "fmt"

const maxInt64 = int64(1<<63 - 1)

type LineItem struct {
	Quantity       int64
	UnitPriceMinor int64
	TaxRateBPS     int64
	DiscountMinor  int64
}

type Totals struct {
	SubtotalMinor int64
	DiscountMinor int64
	TaxMinor      int64
	TotalMinor    int64
}

func CalculateTotals(items []LineItem) (Totals, error) {
	var totals Totals
	for index, item := range items {
		if item.Quantity <= 0 {
			return Totals{}, fmt.Errorf("line item %d: quantity must be greater than zero", index)
		}
		if item.UnitPriceMinor < 0 {
			return Totals{}, fmt.Errorf("line item %d: unit price cannot be negative", index)
		}
		if item.TaxRateBPS < 0 || item.TaxRateBPS > 10000 {
			return Totals{}, fmt.Errorf("line item %d: tax rate must be between 0 and 10000 basis points", index)
		}
		if item.DiscountMinor < 0 {
			return Totals{}, fmt.Errorf("line item %d: discount cannot be negative", index)
		}

		lineSubtotal, err := checkedMul(item.Quantity, item.UnitPriceMinor)
		if err != nil {
			return Totals{}, fmt.Errorf("line item %d subtotal: %w", index, err)
		}
		if item.DiscountMinor > lineSubtotal {
			return Totals{}, fmt.Errorf("line item %d: discount exceeds line subtotal", index)
		}
		taxableMinor := lineSubtotal - item.DiscountMinor
		taxBase, err := checkedMul(taxableMinor, item.TaxRateBPS)
		if err != nil {
			return Totals{}, fmt.Errorf("line item %d tax: %w", index, err)
		}
		taxNumerator, err := checkedAdd(taxBase, 5000)
		if err != nil {
			return Totals{}, fmt.Errorf("line item %d tax rounding: %w", index, err)
		}
		lineTax := taxNumerator / 10000

		totals.SubtotalMinor, err = checkedAdd(totals.SubtotalMinor, lineSubtotal)
		if err != nil {
			return Totals{}, fmt.Errorf("subtotal: %w", err)
		}
		totals.DiscountMinor, err = checkedAdd(totals.DiscountMinor, item.DiscountMinor)
		if err != nil {
			return Totals{}, fmt.Errorf("discount: %w", err)
		}
		totals.TaxMinor, err = checkedAdd(totals.TaxMinor, lineTax)
		if err != nil {
			return Totals{}, fmt.Errorf("tax: %w", err)
		}
	}

	netMinor := totals.SubtotalMinor - totals.DiscountMinor
	totalMinor, err := checkedAdd(netMinor, totals.TaxMinor)
	if err != nil {
		return Totals{}, fmt.Errorf("total: %w", err)
	}
	totals.TotalMinor = totalMinor
	return totals, nil
}

func checkedAdd(left, right int64) (int64, error) {
	if right > 0 && left > maxInt64-right {
		return 0, fmt.Errorf("integer overflow")
	}
	return left + right, nil
}

func checkedMul(left, right int64) (int64, error) {
	if left < 0 || right < 0 {
		return 0, fmt.Errorf("negative arithmetic operand")
	}
	if left != 0 && right > maxInt64/left {
		return 0, fmt.Errorf("integer overflow")
	}
	return left * right, nil
}
