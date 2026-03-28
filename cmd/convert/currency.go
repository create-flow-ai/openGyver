package convert

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Supported currency codes — all recognized aliases map to ISO 4217 codes.
var currencyCodes = map[string]string{
	"usd": "USD", "eur": "EUR", "gbp": "GBP", "jpy": "JPY",
	"cad": "CAD", "aud": "AUD", "chf": "CHF", "cny": "CNY",
	"inr": "INR", "mxn": "MXN", "brl": "BRL", "krw": "KRW",
	"sgd": "SGD", "hkd": "HKD", "nok": "NOK", "sek": "SEK",
	"dkk": "DKK", "nzd": "NZD", "zar": "ZAR", "rub": "RUB",
	"try": "TRY", "pln": "PLN", "thb": "THB", "idr": "IDR",
	"huf": "HUF", "czk": "CZK", "ils": "ILS", "clp": "CLP",
	"php": "PHP", "aed": "AED", "cop": "COP", "sar": "SAR",
	"myr": "MYR", "ron": "RON", "bgn": "BGN", "hrk": "HRK",
	"isk": "ISK", "twd": "TWD",
}

func init() {
	// Build Unit map from currency codes so lookup() can find them.
	units := make(map[string]Unit, len(currencyCodes))
	for alias, code := range currencyCodes {
		units[alias] = Unit{Name: code}
	}

	registerCategory(Category{
		Name:  "Currency",
		Base:  "usd",
		Units: units,
		Convert: func(val float64, from, to string) (float64, error) {
			fromCode := currencyCodes[strings.ToLower(from)]
			toCode := currencyCodes[strings.ToLower(to)]
			if fromCode == "" || toCode == "" {
				return 0, fmt.Errorf("unknown currency pair: %s → %s", from, to)
			}
			return fetchRate(val, fromCode, toCode)
		},
	})
}

// frankfurterResponse is the JSON shape returned by the Frankfurter API.
type frankfurterResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// fetchRate calls the Frankfurter API (free, no key required) for live exchange rates.
func fetchRate(amount float64, from, to string) (float64, error) {
	url := fmt.Sprintf("https://api.frankfurter.app/latest?amount=%.6f&from=%s&to=%s",
		amount, from, to)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("currency API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("currency API returned status %d", resp.StatusCode)
	}

	var result frankfurterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse currency API response: %w", err)
	}

	rate, ok := result.Rates[to]
	if !ok {
		return 0, fmt.Errorf("no rate returned for %s", to)
	}
	return rate, nil
}
