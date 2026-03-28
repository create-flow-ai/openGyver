package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	abbreviated bool
)

var convertCmd = &cobra.Command{
	Use:   "convert <value> <from-unit> <to-unit>",
	Short: "Convert between units of measurement",
	Long: `Convert values between units. Automatically detects the category from the unit names.

CATEGORIES AND UNITS:

  Temperature     c, f, k (celsius, fahrenheit, kelvin)
  Length          mm, cm, m, km, in, ft, yd, mi, nm
  Weight          mg, g, kg, oz, lb, ton, tonne, st (stone)
  Volume          ml, l, gal, qt, pt, cup, floz, tbsp, tsp
  Area            sqmm, sqcm, sqm, sqkm, sqin, sqft, sqyd, sqmi, acre, hectare
  Speed           mps, kph, mph, knots, fps
  Data            b, kb, mb, gb, tb, pb, bit, kbit, mbit, gbit
  Time            ms, sec, min, hr, day, week, month, year
  Currency        usd, eur, gbp, jpy, cad, aud, chf, cny, inr, mxn, brl,
                  krw, sgd, hkd, nok, sek, dkk, nzd, zar, rub, try, pln,
                  thb, idr, huf, czk, ils, clp, php, aed, cop, sar, myr,
                  ron, bgn, hrk, isk, twd
                  (live rates via Frankfurter API — no API key needed)

EXAMPLES:

  openGyver convert 100 cm in          # length
  openGyver convert 72 f c             # temperature
  openGyver convert 1.5 gb mb          # data
  openGyver convert 365 days hours     # time
  openGyver convert 60 mph kph         # speed
  openGyver convert 2.5 acre sqft      # area
  openGyver convert 500 ml cup         # volume
  openGyver convert 150 lb kg          # weight
  openGyver convert 100 usd eur        # currency (live)

Unit names are case-insensitive. Both short and long forms work (e.g. "cm" or "centimeter").`,
	Args: cobra.ExactArgs(3),
	RunE: runConvert,
}

func runConvert(c *cobra.Command, args []string) error {
	valStr, from, to := args[0], strings.ToLower(args[1]), strings.ToLower(args[2])

	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fmt.Errorf("invalid number: %s", valStr)
	}

	fromCat, fromOk := lookup(from)
	toCat, toOk := lookup(to)

	if !fromOk {
		return fmt.Errorf("unknown unit: %s\nRun 'openGyver convert --help' to see supported units", from)
	}
	if !toOk {
		return fmt.Errorf("unknown unit: %s\nRun 'openGyver convert --help' to see supported units", to)
	}
	if fromCat.Name != toCat.Name {
		return fmt.Errorf("cannot convert between %s (%s) and %s (%s)", from, fromCat.Name, to, toCat.Name)
	}

	var result float64
	if fromCat.Convert != nil {
		result, err = fromCat.Convert(val, from, to)
		if err != nil {
			return err
		}
	} else {
		fromUnit := fromCat.Units[from]
		toUnit := fromCat.Units[to]
		result = convertByFactor(val, fromUnit, toUnit)
	}

	toName := fromCat.Units[to].Name
	if abbreviated {
		fmt.Printf("%s %s\n", formatNumber(result), toName)
	} else {
		fromName := fromCat.Units[from].Name
		fmt.Printf("%s %s = %s %s\n",
			formatNumber(val), fromName,
			formatNumber(result), toName)
	}

	return nil
}

func formatNumber(v float64) string {
	if v == float64(int64(v)) && v < 1e15 {
		return fmt.Sprintf("%d", int64(v))
	}
	s := fmt.Sprintf("%.6f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func init() {
	convertCmd.Flags().BoolVarP(&abbreviated, "abbreviated", "a", false, "output only the converted value and unit")
	cmd.Register(convertCmd)
}
