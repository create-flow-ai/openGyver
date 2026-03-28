package finance

import (
	"fmt"
	"math"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

// jsonOut controls JSON output across all subcommands.
var jsonOut bool

// ─── Parent command ─────────────────────────────────────────────────────────

var financeCmd = &cobra.Command{
	Use:   "finance",
	Short: "Financial calculators and converters",
	Long: `A collection of financial calculators for everyday money math.

SUBCOMMANDS:

  loan       Loan/mortgage payment calculator
  compound   Compound interest calculator
  roi        Return on investment calculator
  tip        Tip calculator with bill splitting
  tax        Sales tax calculator
  salary     Salary converter between pay periods
  discount   Discount/sale price calculator
  margin     Profit margin and markup calculator

All subcommands support --json / -j for machine-readable output.

EXAMPLES:

  openGyver finance loan --principal 250000 --rate 6.5 --years 30
  openGyver finance compound --principal 10000 --rate 7 --years 10
  openGyver finance roi --initial 1000 --final 1500
  openGyver finance tip --amount 85.50 --percent 20 --split 4
  openGyver finance tax --amount 99.99 --rate 8.25
  openGyver finance salary --amount 50 --from hourly --to yearly
  openGyver finance discount --price 199.99 --percent 25
  openGyver finance margin --cost 40 --revenue 100`,
}

func init() {
	financeCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	financeCmd.AddCommand(loanCmd)
	financeCmd.AddCommand(compoundCmd)
	financeCmd.AddCommand(roiCmd)
	financeCmd.AddCommand(tipCmd)
	financeCmd.AddCommand(taxCmd)
	financeCmd.AddCommand(salaryCmd)
	financeCmd.AddCommand(discountCmd)
	financeCmd.AddCommand(marginCmd)
	cmd.Register(financeCmd)
}

// ─── loan subcommand ───────────────────────────────────────────────────────

var (
	loanPrincipal float64
	loanRate      float64
	loanYears     int
)

var loanCmd = &cobra.Command{
	Use:   "loan",
	Short: "Loan/mortgage payment calculator",
	Long: `Calculate monthly payment, total payment, and total interest for a
fixed-rate loan or mortgage.

Uses the standard amortisation formula:

  M = P * [r(1+r)^n] / [(1+r)^n - 1]

where P = principal, r = monthly interest rate, n = total number of
monthly payments.

FLAGS:

  --principal   Loan amount in dollars (required)
  --rate        Annual interest rate as a percentage (required)
  --years       Loan term in years (required)

EXAMPLES:

  # 30-year mortgage at 6.5%
  openGyver finance loan --principal 250000 --rate 6.5 --years 30

  # 5-year car loan at 4.9%
  openGyver finance loan --principal 35000 --rate 4.9 --years 5

  # 15-year mortgage at 5.75% (JSON output)
  openGyver finance loan --principal 400000 --rate 5.75 --years 15 -j`,
	RunE: runLoan,
}

func init() {
	loanCmd.Flags().Float64Var(&loanPrincipal, "principal", 0, "loan amount in dollars (required)")
	loanCmd.Flags().Float64Var(&loanRate, "rate", 0, "annual interest rate as a percentage (required)")
	loanCmd.Flags().IntVar(&loanYears, "years", 0, "loan term in years (required)")
	_ = loanCmd.MarkFlagRequired("principal")
	_ = loanCmd.MarkFlagRequired("rate")
	_ = loanCmd.MarkFlagRequired("years")
}

func runLoan(_ *cobra.Command, _ []string) error {
	if loanPrincipal <= 0 {
		return fmt.Errorf("--principal must be greater than 0")
	}
	if loanRate <= 0 {
		return fmt.Errorf("--rate must be greater than 0")
	}
	if loanYears <= 0 {
		return fmt.Errorf("--years must be greater than 0")
	}

	r := loanRate / 100.0 / 12.0 // monthly interest rate
	n := float64(loanYears * 12)  // total number of payments
	P := loanPrincipal

	// M = P * [r(1+r)^n] / [(1+r)^n - 1]
	pow := math.Pow(1+r, n)
	monthly := P * (r * pow) / (pow - 1)
	totalPayment := monthly * n
	totalInterest := totalPayment - P

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"principal":      P,
			"annual_rate":    loanRate,
			"years":          loanYears,
			"monthly_payment": round2(monthly),
			"total_payment":   round2(totalPayment),
			"total_interest":  round2(totalInterest),
		})
	}

	fmt.Printf("Loan Calculator\n")
	fmt.Printf("  Principal:        $%s\n", formatMoney(P))
	fmt.Printf("  Annual Rate:      %.2f%%\n", loanRate)
	fmt.Printf("  Term:             %d years (%d payments)\n", loanYears, loanYears*12)
	fmt.Println()
	fmt.Printf("  Monthly Payment:  $%s\n", formatMoney(monthly))
	fmt.Printf("  Total Payment:    $%s\n", formatMoney(totalPayment))
	fmt.Printf("  Total Interest:   $%s\n", formatMoney(totalInterest))
	return nil
}

// ─── compound subcommand ───────────────────────────────────────────────────

var (
	compPrincipal float64
	compRate      float64
	compYears     int
	compFrequency int
)

var compoundCmd = &cobra.Command{
	Use:   "compound",
	Short: "Compound interest calculator",
	Long: `Calculate the future value of an investment with compound interest.

Uses the compound interest formula:

  A = P * (1 + r/n)^(n*t)

where P = principal, r = annual rate, n = compounding frequency per
year, t = time in years.

FLAGS:

  --principal    Initial investment amount (required)
  --rate         Annual interest rate as a percentage (required)
  --years        Investment period in years (required)
  --frequency    Compounding frequency per year (default 12 = monthly)

COMMON FREQUENCIES:

  1    = annually
  4    = quarterly
  12   = monthly (default)
  52   = weekly
  365  = daily

EXAMPLES:

  # $10,000 at 7% for 10 years, compounded monthly
  openGyver finance compound --principal 10000 --rate 7 --years 10

  # $5,000 at 5% for 20 years, compounded daily
  openGyver finance compound --principal 5000 --rate 5 --years 20 --frequency 365

  # $25,000 at 8.5% for 5 years, compounded quarterly
  openGyver finance compound --principal 25000 --rate 8.5 --years 5 --frequency 4`,
	RunE: runCompound,
}

func init() {
	compoundCmd.Flags().Float64Var(&compPrincipal, "principal", 0, "initial investment amount (required)")
	compoundCmd.Flags().Float64Var(&compRate, "rate", 0, "annual interest rate as a percentage (required)")
	compoundCmd.Flags().IntVar(&compYears, "years", 0, "investment period in years (required)")
	compoundCmd.Flags().IntVar(&compFrequency, "frequency", 12, "compounding frequency per year (default 12 = monthly)")
	_ = compoundCmd.MarkFlagRequired("principal")
	_ = compoundCmd.MarkFlagRequired("rate")
	_ = compoundCmd.MarkFlagRequired("years")
}

func runCompound(_ *cobra.Command, _ []string) error {
	if compPrincipal <= 0 {
		return fmt.Errorf("--principal must be greater than 0")
	}
	if compRate <= 0 {
		return fmt.Errorf("--rate must be greater than 0")
	}
	if compYears <= 0 {
		return fmt.Errorf("--years must be greater than 0")
	}
	if compFrequency <= 0 {
		return fmt.Errorf("--frequency must be greater than 0")
	}

	P := compPrincipal
	r := compRate / 100.0
	n := float64(compFrequency)
	t := float64(compYears)

	// A = P * (1 + r/n)^(n*t)
	finalAmount := P * math.Pow(1+r/n, n*t)
	totalInterest := finalAmount - P

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"principal":      P,
			"annual_rate":    compRate,
			"years":          compYears,
			"frequency":      compFrequency,
			"final_amount":   round2(finalAmount),
			"total_interest": round2(totalInterest),
		})
	}

	fmt.Printf("Compound Interest Calculator\n")
	fmt.Printf("  Principal:        $%s\n", formatMoney(P))
	fmt.Printf("  Annual Rate:      %.2f%%\n", compRate)
	fmt.Printf("  Term:             %d years\n", compYears)
	fmt.Printf("  Compounding:      %s\n", frequencyLabel(compFrequency))
	fmt.Println()
	fmt.Printf("  Final Amount:     $%s\n", formatMoney(finalAmount))
	fmt.Printf("  Total Interest:   $%s\n", formatMoney(totalInterest))
	return nil
}

// ─── roi subcommand ────────────────────────────────────────────────────────

var (
	roiInitial float64
	roiFinal   float64
)

var roiCmd = &cobra.Command{
	Use:   "roi",
	Short: "Return on investment calculator",
	Long: `Calculate the return on investment (ROI) given an initial investment
and a final value.

  ROI = ((final - initial) / initial) * 100

Positive ROI indicates a profit; negative ROI indicates a loss.

FLAGS:

  --initial   Initial investment amount (required)
  --final     Final value of the investment (required)

EXAMPLES:

  # Bought at $1,000, now worth $1,500
  openGyver finance roi --initial 1000 --final 1500

  # Bought at $5,000, sold for $3,750 (a loss)
  openGyver finance roi --initial 5000 --final 3750

  # JSON output
  openGyver finance roi --initial 200 --final 350 -j`,
	RunE: runROI,
}

func init() {
	roiCmd.Flags().Float64Var(&roiInitial, "initial", 0, "initial investment amount (required)")
	roiCmd.Flags().Float64Var(&roiFinal, "final", 0, "final value of the investment (required)")
	_ = roiCmd.MarkFlagRequired("initial")
	_ = roiCmd.MarkFlagRequired("final")
}

func runROI(_ *cobra.Command, _ []string) error {
	if roiInitial == 0 {
		return fmt.Errorf("--initial must not be zero")
	}

	profitLoss := roiFinal - roiInitial
	roiPercent := (profitLoss / roiInitial) * 100.0

	label := "Profit"
	if profitLoss < 0 {
		label = "Loss"
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"initial":     roiInitial,
			"final":       roiFinal,
			"profit_loss": round2(profitLoss),
			"roi_percent": round2(roiPercent),
		})
	}

	fmt.Printf("Return on Investment\n")
	fmt.Printf("  Initial Investment: $%s\n", formatMoney(roiInitial))
	fmt.Printf("  Final Value:        $%s\n", formatMoney(roiFinal))
	fmt.Println()
	fmt.Printf("  %s:  %s$%s\n", label, sign(profitLoss), formatMoney(math.Abs(profitLoss)))
	fmt.Printf("  ROI:                %s%.2f%%\n", sign(roiPercent), math.Abs(roiPercent))
	return nil
}

// ─── tip subcommand ────────────────────────────────────────────────────────

var (
	tipAmount  float64
	tipPercent float64
	tipSplit   int
)

var tipCmd = &cobra.Command{
	Use:   "tip",
	Short: "Tip calculator with bill splitting",
	Long: `Calculate tip, total bill, and per-person share when splitting
the bill among multiple people.

FLAGS:

  --amount    Bill amount before tip (required)
  --percent   Tip percentage (default 18)
  --split     Number of people splitting the bill (default 1)

EXAMPLES:

  # Standard 18% tip on $85.50
  openGyver finance tip --amount 85.50

  # 20% tip on $120, split 4 ways
  openGyver finance tip --amount 120 --percent 20 --split 4

  # 15% tip on $45.00
  openGyver finance tip --amount 45 --percent 15`,
	RunE: runTip,
}

func init() {
	tipCmd.Flags().Float64Var(&tipAmount, "amount", 0, "bill amount before tip (required)")
	tipCmd.Flags().Float64Var(&tipPercent, "percent", 18, "tip percentage (default 18)")
	tipCmd.Flags().IntVar(&tipSplit, "split", 1, "number of people splitting the bill (default 1)")
	_ = tipCmd.MarkFlagRequired("amount")
}

func runTip(_ *cobra.Command, _ []string) error {
	if tipAmount <= 0 {
		return fmt.Errorf("--amount must be greater than 0")
	}
	if tipPercent < 0 {
		return fmt.Errorf("--percent must not be negative")
	}
	if tipSplit <= 0 {
		return fmt.Errorf("--split must be greater than 0")
	}

	tip := tipAmount * tipPercent / 100.0
	total := tipAmount + tip
	perPerson := total / float64(tipSplit)

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"amount":     tipAmount,
			"percent":    tipPercent,
			"tip":        round2(tip),
			"total":      round2(total),
			"split":      tipSplit,
			"per_person": round2(perPerson),
		})
	}

	fmt.Printf("Tip Calculator\n")
	fmt.Printf("  Bill Amount:   $%s\n", formatMoney(tipAmount))
	fmt.Printf("  Tip (%.0f%%):     $%s\n", tipPercent, formatMoney(tip))
	fmt.Printf("  Total:         $%s\n", formatMoney(total))
	if tipSplit > 1 {
		fmt.Printf("  Split %d ways:  $%s per person\n", tipSplit, formatMoney(perPerson))
	}
	return nil
}

// ─── tax subcommand ────────────────────────────────────────────────────────

var (
	taxAmount float64
	taxRate   float64
)

var taxCmd = &cobra.Command{
	Use:   "tax",
	Short: "Sales tax calculator",
	Long: `Calculate the tax amount and total price given a pre-tax amount
and a tax rate.

FLAGS:

  --amount   Pre-tax amount (required)
  --rate     Tax rate as a percentage (required)

EXAMPLES:

  # 8.25% sales tax on $99.99
  openGyver finance tax --amount 99.99 --rate 8.25

  # 10% tax on $250
  openGyver finance tax --amount 250 --rate 10

  # JSON output
  openGyver finance tax --amount 49.95 --rate 7.5 -j`,
	RunE: runTax,
}

func init() {
	taxCmd.Flags().Float64Var(&taxAmount, "amount", 0, "pre-tax amount (required)")
	taxCmd.Flags().Float64Var(&taxRate, "rate", 0, "tax rate as a percentage (required)")
	_ = taxCmd.MarkFlagRequired("amount")
	_ = taxCmd.MarkFlagRequired("rate")
}

func runTax(_ *cobra.Command, _ []string) error {
	if taxAmount <= 0 {
		return fmt.Errorf("--amount must be greater than 0")
	}
	if taxRate < 0 {
		return fmt.Errorf("--rate must not be negative")
	}

	tax := taxAmount * taxRate / 100.0
	total := taxAmount + tax

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"amount": taxAmount,
			"rate":   taxRate,
			"tax":    round2(tax),
			"total":  round2(total),
		})
	}

	fmt.Printf("Tax Calculator\n")
	fmt.Printf("  Amount:   $%s\n", formatMoney(taxAmount))
	fmt.Printf("  Tax Rate: %.2f%%\n", taxRate)
	fmt.Println()
	fmt.Printf("  Tax:      $%s\n", formatMoney(tax))
	fmt.Printf("  Total:    $%s\n", formatMoney(total))
	return nil
}

// ─── salary subcommand ─────────────────────────────────────────────────────

var (
	salaryAmount float64
	salaryFrom   string
	salaryTo     string
)

var salaryCmd = &cobra.Command{
	Use:   "salary",
	Short: "Salary converter between pay periods",
	Long: `Convert salary or wages between different pay periods.

Assumes standard full-time work schedule:
  - 8 hours per day
  - 5 days per week
  - 4.33 weeks per month (52/12)
  - 52 weeks per year

Valid period values: hourly, daily, weekly, monthly, yearly

FLAGS:

  --amount   Salary/wage amount (required)
  --from     Source pay period (required)
  --to       Target pay period (required)

EXAMPLES:

  # Hourly to yearly
  openGyver finance salary --amount 50 --from hourly --to yearly

  # Yearly to hourly
  openGyver finance salary --amount 100000 --from yearly --to hourly

  # Monthly to weekly
  openGyver finance salary --amount 8000 --from monthly --to weekly

  # Daily to monthly
  openGyver finance salary --amount 400 --from daily --to monthly`,
	RunE: runSalary,
}

func init() {
	salaryCmd.Flags().Float64Var(&salaryAmount, "amount", 0, "salary/wage amount (required)")
	salaryCmd.Flags().StringVar(&salaryFrom, "from", "", "source pay period: hourly, daily, weekly, monthly, yearly (required)")
	salaryCmd.Flags().StringVar(&salaryTo, "to", "", "target pay period: hourly, daily, weekly, monthly, yearly (required)")
	_ = salaryCmd.MarkFlagRequired("amount")
	_ = salaryCmd.MarkFlagRequired("from")
	_ = salaryCmd.MarkFlagRequired("to")
}

// salaryToHourly maps each pay period to the multiplier that converts it to hourly.
// Assumptions: 8h/day, 5d/week, 4.33w/month (52/12), 52w/year.
var salaryToHourly = map[string]float64{
	"hourly":  1,
	"daily":   1.0 / 8.0,
	"weekly":  1.0 / (8.0 * 5.0),
	"monthly": 1.0 / (8.0 * 5.0 * (52.0 / 12.0)),
	"yearly":  1.0 / (8.0 * 5.0 * 52.0),
}

// hourlyToSalary maps each pay period to the multiplier that converts hourly to it.
var hourlyToSalary = map[string]float64{
	"hourly":  1,
	"daily":   8.0,
	"weekly":  8.0 * 5.0,
	"monthly": 8.0 * 5.0 * (52.0 / 12.0),
	"yearly":  8.0 * 5.0 * 52.0,
}

var validPeriods = []string{"hourly", "daily", "weekly", "monthly", "yearly"}

func runSalary(_ *cobra.Command, _ []string) error {
	from := strings.ToLower(salaryFrom)
	to := strings.ToLower(salaryTo)

	toHourly, ok := salaryToHourly[from]
	if !ok {
		return fmt.Errorf("invalid --from period %q; valid values: %s", salaryFrom, strings.Join(validPeriods, ", "))
	}
	fromHourly, ok := hourlyToSalary[to]
	if !ok {
		return fmt.Errorf("invalid --to period %q; valid values: %s", salaryTo, strings.Join(validPeriods, ", "))
	}

	if salaryAmount <= 0 {
		return fmt.Errorf("--amount must be greater than 0")
	}

	// Convert source amount -> hourly -> target
	hourly := salaryAmount * toHourly
	result := hourly * fromHourly

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"amount": salaryAmount,
			"from":   from,
			"to":     to,
			"result": round2(result),
		})
	}

	fmt.Printf("Salary Converter\n")
	fmt.Printf("  $%s %s\n", formatMoney(salaryAmount), from)
	fmt.Println()
	fmt.Printf("  Hourly:   $%s\n", formatMoney(hourly))
	fmt.Printf("  Daily:    $%s\n", formatMoney(hourly*8))
	fmt.Printf("  Weekly:   $%s\n", formatMoney(hourly*8*5))
	fmt.Printf("  Monthly:  $%s\n", formatMoney(hourly*8*5*(52.0/12.0)))
	fmt.Printf("  Yearly:   $%s\n", formatMoney(hourly*8*5*52))
	return nil
}

// ─── discount subcommand ───────────────────────────────────────────────────

var (
	discountPrice   float64
	discountPercent float64
)

var discountCmd = &cobra.Command{
	Use:   "discount",
	Short: "Discount/sale price calculator",
	Long: `Calculate the discount amount and final price after applying a
percentage discount.

FLAGS:

  --price     Original price (required)
  --percent   Discount percentage (required)

EXAMPLES:

  # 25% off $199.99
  openGyver finance discount --price 199.99 --percent 25

  # 50% off $80
  openGyver finance discount --price 80 --percent 50

  # 10% off $1,250
  openGyver finance discount --price 1250 --percent 10 -j`,
	RunE: runDiscount,
}

func init() {
	discountCmd.Flags().Float64Var(&discountPrice, "price", 0, "original price (required)")
	discountCmd.Flags().Float64Var(&discountPercent, "percent", 0, "discount percentage (required)")
	_ = discountCmd.MarkFlagRequired("price")
	_ = discountCmd.MarkFlagRequired("percent")
}

func runDiscount(_ *cobra.Command, _ []string) error {
	if discountPrice <= 0 {
		return fmt.Errorf("--price must be greater than 0")
	}
	if discountPercent < 0 || discountPercent > 100 {
		return fmt.Errorf("--percent must be between 0 and 100")
	}

	discountAmount := discountPrice * discountPercent / 100.0
	finalPrice := discountPrice - discountAmount

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"original_price":  discountPrice,
			"discount_percent": discountPercent,
			"discount_amount":  round2(discountAmount),
			"final_price":      round2(finalPrice),
		})
	}

	fmt.Printf("Discount Calculator\n")
	fmt.Printf("  Original Price:   $%s\n", formatMoney(discountPrice))
	fmt.Printf("  Discount:         %.0f%%\n", discountPercent)
	fmt.Println()
	fmt.Printf("  You Save:         $%s\n", formatMoney(discountAmount))
	fmt.Printf("  Final Price:      $%s\n", formatMoney(finalPrice))
	return nil
}

// ─── margin subcommand ─────────────────────────────────────────────────────

var (
	marginCost    float64
	marginRevenue float64
)

var marginCmd = &cobra.Command{
	Use:   "margin",
	Short: "Profit margin and markup calculator",
	Long: `Calculate profit, profit margin percentage, and markup percentage
given cost and revenue (selling price).

  Profit     = Revenue - Cost
  Margin %   = (Profit / Revenue) * 100
  Markup %   = (Profit / Cost) * 100

FLAGS:

  --cost      Cost of goods/services (required)
  --revenue   Selling price / revenue (required)

EXAMPLES:

  # Cost $40, selling for $100
  openGyver finance margin --cost 40 --revenue 100

  # Cost $15.50, selling for $29.99
  openGyver finance margin --cost 15.50 --revenue 29.99

  # JSON output
  openGyver finance margin --cost 250 --revenue 400 -j`,
	RunE: runMargin,
}

func init() {
	marginCmd.Flags().Float64Var(&marginCost, "cost", 0, "cost of goods/services (required)")
	marginCmd.Flags().Float64Var(&marginRevenue, "revenue", 0, "selling price / revenue (required)")
	_ = marginCmd.MarkFlagRequired("cost")
	_ = marginCmd.MarkFlagRequired("revenue")
}

func runMargin(_ *cobra.Command, _ []string) error {
	if marginCost <= 0 {
		return fmt.Errorf("--cost must be greater than 0")
	}
	if marginRevenue <= 0 {
		return fmt.Errorf("--revenue must be greater than 0")
	}

	profit := marginRevenue - marginCost
	marginPercent := (profit / marginRevenue) * 100.0
	markupPercent := (profit / marginCost) * 100.0

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"cost":           marginCost,
			"revenue":        marginRevenue,
			"profit":         round2(profit),
			"margin_percent": round2(marginPercent),
			"markup_percent": round2(markupPercent),
		})
	}

	fmt.Printf("Profit Margin Calculator\n")
	fmt.Printf("  Cost:      $%s\n", formatMoney(marginCost))
	fmt.Printf("  Revenue:   $%s\n", formatMoney(marginRevenue))
	fmt.Println()
	fmt.Printf("  Profit:    %s$%s\n", sign(profit), formatMoney(math.Abs(profit)))
	fmt.Printf("  Margin:    %.2f%%\n", marginPercent)
	fmt.Printf("  Markup:    %.2f%%\n", markupPercent)
	return nil
}

// ─── Helpers ───────────────────────────────────────────────────────────────

// round2 rounds a float to 2 decimal places.
func round2(f float64) float64 {
	return math.Round(f*100) / 100
}

// formatMoney formats a float as a dollar amount with commas and 2 decimal
// places (e.g., 1,234,567.89).
func formatMoney(f float64) string {
	negative := f < 0
	if negative {
		f = -f
	}

	// Format with 2 decimal places.
	s := fmt.Sprintf("%.2f", f)

	// Split into integer and decimal parts.
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	decPart := parts[1]

	// Insert commas into the integer part.
	n := len(intPart)
	if n > 3 {
		var b strings.Builder
		remainder := n % 3
		if remainder > 0 {
			b.WriteString(intPart[:remainder])
		}
		for i := remainder; i < n; i += 3 {
			if b.Len() > 0 {
				b.WriteByte(',')
			}
			b.WriteString(intPart[i : i+3])
		}
		intPart = b.String()
	}

	result := intPart + "." + decPart
	if negative {
		return "-" + result
	}
	return result
}

// sign returns "+" for non-negative values and "-" for negative values.
func sign(f float64) string {
	if f < 0 {
		return "-"
	}
	return ""
}

// frequencyLabel returns a human-readable label for a compounding frequency.
func frequencyLabel(freq int) string {
	switch freq {
	case 1:
		return "annually (1x/year)"
	case 2:
		return "semi-annually (2x/year)"
	case 4:
		return "quarterly (4x/year)"
	case 12:
		return "monthly (12x/year)"
	case 26:
		return "bi-weekly (26x/year)"
	case 52:
		return "weekly (52x/year)"
	case 365:
		return "daily (365x/year)"
	default:
		return fmt.Sprintf("%dx/year", freq)
	}
}
