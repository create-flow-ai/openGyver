package finance

import (
	"math"
	"testing"
)

func TestRound2(t *testing.T) {
	tests := []struct {
		in   float64
		want float64
	}{
		{1.234, 1.23},
		{1.235, 1.24},
		{1.999, 2.00},
		{0.001, 0.00},
	}
	for _, tt := range tests {
		got := round2(tt.in)
		if got != tt.want {
			t.Errorf("round2(%f) = %f, want %f", tt.in, got, tt.want)
		}
	}
}

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		in   float64
		want string
	}{
		{1234567.89, "1,234,567.89"},
		{0.50, "0.50"},
		{999.99, "999.99"},
		{1000.00, "1,000.00"},
	}
	for _, tt := range tests {
		got := formatMoney(tt.in)
		if got != tt.want {
			t.Errorf("formatMoney(%f) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestLoanCalculation(t *testing.T) {
	// Known values: $300,000 principal, 6.5% annual rate, 30 years.
	// Expected monthly payment: ~$1896.20
	P := 300000.0
	rate := 6.5
	years := 30

	r := rate / 100.0 / 12.0
	n := float64(years * 12)
	pow := math.Pow(1+r, n)
	monthly := P * (r * pow) / (pow - 1)

	expected := 1896.20
	if math.Abs(monthly-expected) > 0.10 {
		t.Errorf("loan monthly payment = %.2f, want ~%.2f", monthly, expected)
	}

	totalPayment := monthly * n
	totalInterest := totalPayment - P
	if totalInterest <= 0 {
		t.Error("total interest should be positive")
	}
	if totalPayment <= P {
		t.Error("total payment should exceed principal")
	}
}

func TestCompoundInterest(t *testing.T) {
	// $10,000 at 7% for 10 years, compounded monthly.
	P := 10000.0
	r := 7.0 / 100.0
	n := 12.0
	tYears := 10.0

	finalAmount := P * math.Pow(1+r/n, n*tYears)
	// Expected approximately $20,096.61
	if finalAmount < 20000 || finalAmount > 20200 {
		t.Errorf("compound interest result = %.2f, expected ~20096", finalAmount)
	}

	totalInterest := finalAmount - P
	if totalInterest <= 0 {
		t.Error("total interest should be positive")
	}
}

func TestROI(t *testing.T) {
	// Initial $1000, final $1500 -> ROI = 50%.
	initial := 1000.0
	final := 1500.0
	profitLoss := final - initial
	roiPercent := (profitLoss / initial) * 100.0

	if round2(roiPercent) != 50.0 {
		t.Errorf("ROI = %.2f%%, want 50.00%%", roiPercent)
	}
	if profitLoss != 500.0 {
		t.Errorf("profit/loss = %.2f, want 500.00", profitLoss)
	}

	// Loss case.
	initial2 := 5000.0
	final2 := 3750.0
	roi2 := ((final2 - initial2) / initial2) * 100.0
	if round2(roi2) != -25.0 {
		t.Errorf("ROI (loss) = %.2f%%, want -25.00%%", roi2)
	}
}

func TestTipCalculation(t *testing.T) {
	amount := 85.50
	percent := 20.0
	split := 4

	tip := amount * percent / 100.0
	total := amount + tip
	perPerson := total / float64(split)

	if round2(tip) != 17.10 {
		t.Errorf("tip = %.2f, want 17.10", tip)
	}
	if round2(total) != 102.60 {
		t.Errorf("total = %.2f, want 102.60", total)
	}
	if round2(perPerson) != 25.65 {
		t.Errorf("perPerson = %.2f, want 25.65", perPerson)
	}
}

func TestSalaryConversion(t *testing.T) {
	// Hourly $50 -> yearly should be 50 * 8 * 5 * 52 = 104,000.
	hourlyRate := 50.0
	toHourly := salaryToHourly["hourly"]   // 1.0
	fromHourly := hourlyToSalary["yearly"] // 8 * 5 * 52 = 2080

	hourly := hourlyRate * toHourly
	yearly := hourly * fromHourly

	if round2(yearly) != 104000.0 {
		t.Errorf("salary hourly->yearly = %.2f, want 104000.00", yearly)
	}

	// Reverse: yearly $104,000 -> hourly should be $50.
	toHourly2 := salaryToHourly["yearly"]
	fromHourly2 := hourlyToSalary["hourly"]
	h := 104000.0 * toHourly2 * fromHourly2
	if round2(h) != 50.0 {
		t.Errorf("salary yearly->hourly = %.2f, want 50.00", h)
	}
}

func TestDiscount(t *testing.T) {
	price := 199.99
	percent := 25.0
	discountAmount := price * percent / 100.0
	finalPrice := price - discountAmount

	if round2(discountAmount) != 50.00 {
		t.Errorf("discount amount = %.2f, want 50.00", discountAmount)
	}
	if round2(finalPrice) != 149.99 {
		t.Errorf("final price = %.2f, want 149.99", finalPrice)
	}
}

func TestMargin(t *testing.T) {
	cost := 40.0
	revenue := 100.0

	profit := revenue - cost
	marginPercent := (profit / revenue) * 100.0
	markupPercent := (profit / cost) * 100.0

	if round2(profit) != 60.0 {
		t.Errorf("profit = %.2f, want 60.00", profit)
	}
	if round2(marginPercent) != 60.0 {
		t.Errorf("margin%% = %.2f, want 60.00", marginPercent)
	}
	if round2(markupPercent) != 150.0 {
		t.Errorf("markup%% = %.2f, want 150.00", markupPercent)
	}
}

func TestSign(t *testing.T) {
	if sign(5.0) != "" {
		t.Errorf("sign(5.0) = %q, want empty string", sign(5.0))
	}
	if sign(-3.0) != "-" {
		t.Errorf("sign(-3.0) = %q, want \"-\"", sign(-3.0))
	}
	if sign(0) != "" {
		t.Errorf("sign(0) = %q, want empty string", sign(0))
	}
}

func TestFrequencyLabel(t *testing.T) {
	if frequencyLabel(12) != "monthly (12x/year)" {
		t.Errorf("frequencyLabel(12) = %q", frequencyLabel(12))
	}
	if frequencyLabel(365) != "daily (365x/year)" {
		t.Errorf("frequencyLabel(365) = %q", frequencyLabel(365))
	}
	if frequencyLabel(1) != "annually (1x/year)" {
		t.Errorf("frequencyLabel(1) = %q", frequencyLabel(1))
	}
}
