package stock

import (
	"testing"

	"github.com/spf13/cobra"
)

// ---------------------------------------------------------------------------
// Command metadata
// ---------------------------------------------------------------------------

func TestStockCmd_Metadata(t *testing.T) {
	if stockCmd.Use != "stock <ticker>" {
		t.Errorf("unexpected Use: %s", stockCmd.Use)
	}
	if stockCmd.Short == "" {
		t.Error("Short should not be empty")
	}
	if stockCmd.Long == "" {
		t.Error("Long should not be empty")
	}
}

func TestStockCmd_RequiresOneArg(t *testing.T) {
	v := cobra.ExactArgs(1)
	if err := v(stockCmd, []string{}); err == nil {
		t.Error("expected error with zero args")
	}
	if err := v(stockCmd, []string{"AAPL"}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := v(stockCmd, []string{"AAPL", "MSFT"}); err == nil {
		t.Error("expected error with two args")
	}
}

func TestStockCmd_Flags(t *testing.T) {
	f := stockCmd.Flags()
	for _, name := range []string{"date", "from", "to", "market", "interval"} {
		if f.Lookup(name) == nil {
			t.Errorf("--%s flag not found", name)
		}
	}
	if f.ShorthandLookup("d") == nil {
		t.Error("-d shorthand not found for --date")
	}
	if f.ShorthandLookup("m") == nil {
		t.Error("-m shorthand not found for --market")
	}
}

// ---------------------------------------------------------------------------
// resolveTicker
// ---------------------------------------------------------------------------

func TestResolveTicker_Plain(t *testing.T) {
	got := resolveTicker("AAPL", "")
	if got != "AAPL" {
		t.Errorf("got %q, want AAPL", got)
	}
}

func TestResolveTicker_Lowercase(t *testing.T) {
	got := resolveTicker("aapl", "")
	if got != "AAPL" {
		t.Errorf("got %q, want AAPL", got)
	}
}

func TestResolveTicker_WithMarket(t *testing.T) {
	tests := map[string]struct {
		ticker, market, want string
	}{
		"kospi":     {"005930", "kospi", "005930.KS"},
		"kosdaq":    {"035720", "kosdaq", "035720.KQ"},
		"tokyo":     {"7203", "tokyo", "7203.T"},
		"london":    {"SHEL", "london", "SHEL.L"},
		"hongkong":  {"0700", "hk", "0700.HK"},
		"shanghai":  {"600519", "shanghai", "600519.SS"},
		"frankfurt": {"SAP", "frankfurt", "SAP.DE"},
		"paris":     {"MC", "paris", "MC.PA"},
		"toronto":   {"RY", "tsx", "RY.TO"},
		"australia": {"BHP", "asx", "BHP.AX"},
		"india_nse": {"RELIANCE", "nse", "RELIANCE.NS"},
		"india_bse": {"RELIANCE", "bse", "RELIANCE.BO"},
		"taiwan":    {"2330", "twse", "2330.TW"},
		"singapore": {"D05", "sgx", "D05.SI"},
		"brazil":    {"PETR4", "bovespa", "PETR4.SA"},
		"swiss":     {"NESN", "six", "NESN.SW"},
	}
	for name, tt := range tests {
		got := resolveTicker(tt.ticker, tt.market)
		if got != tt.want {
			t.Errorf("%s: resolveTicker(%q, %q) = %q, want %q", name, tt.ticker, tt.market, got, tt.want)
		}
	}
}

func TestResolveTicker_AlreadyHasSuffix(t *testing.T) {
	// If ticker already has a dot suffix, don't double-add
	got := resolveTicker("005930.KS", "kospi")
	if got != "005930.KS" {
		t.Errorf("got %q, want 005930.KS (should not double-suffix)", got)
	}
}

func TestResolveTicker_UnknownMarket(t *testing.T) {
	// Unknown market should leave ticker unchanged
	got := resolveTicker("AAPL", "fakexchange")
	if got != "AAPL" {
		t.Errorf("got %q, want AAPL", got)
	}
}

// ---------------------------------------------------------------------------
// Market suffix coverage
// ---------------------------------------------------------------------------

func TestMarketSuffixes_AllNonEmpty(t *testing.T) {
	usMarkets := map[string]bool{"nasdaq": true, "nyse": true, "us": true}
	for name, suffix := range marketSuffixes {
		if !usMarkets[name] && suffix == "" {
			t.Errorf("market %q has empty suffix (should only be empty for US)", name)
		}
	}
}

func TestMarketSuffixes_Count(t *testing.T) {
	if len(marketSuffixes) < 30 {
		t.Errorf("expected at least 30 market mappings, got %d", len(marketSuffixes))
	}
}

// ---------------------------------------------------------------------------
// jsonFloat
// ---------------------------------------------------------------------------

func TestJsonFloat_Null(t *testing.T) {
	var f jsonFloat
	err := f.UnmarshalJSON([]byte("null"))
	if err != nil {
		t.Fatal(err)
	}
	if f != 0 {
		t.Errorf("null should unmarshal to 0, got %f", f)
	}
}

func TestJsonFloat_Number(t *testing.T) {
	var f jsonFloat
	err := f.UnmarshalJSON([]byte("123.45"))
	if err != nil {
		t.Fatal(err)
	}
	if f != 123.45 {
		t.Errorf("got %f, want 123.45", f)
	}
}

func TestJsonFloat_Invalid(t *testing.T) {
	var f jsonFloat
	err := f.UnmarshalJSON([]byte("\"not a number\""))
	if err == nil {
		t.Error("expected error for invalid number")
	}
}

// ---------------------------------------------------------------------------
// Live API tests (run only if network available)
// ---------------------------------------------------------------------------

func TestLookupCurrent_AAPL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}
	err := lookupCurrent("AAPL")
	if err != nil {
		t.Fatalf("lookupCurrent(AAPL) failed: %v", err)
	}
}

func TestLookupDate_AAPL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}
	err := lookupDate("AAPL", "2024-12-20")
	if err != nil {
		t.Fatalf("lookupDate(AAPL, 2024-12-20) failed: %v", err)
	}
}

func TestLookupRange_MSFT(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}
	err := lookupRange("MSFT", "2025-03-01", "2025-03-07", "1d")
	if err != nil {
		t.Fatalf("lookupRange(MSFT) failed: %v", err)
	}
}

func TestLookupCurrent_KoreanStock(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}
	ticker := resolveTicker("005930", "kospi")
	err := lookupCurrent(ticker)
	if err != nil {
		t.Fatalf("lookupCurrent(%s) failed: %v", ticker, err)
	}
}

func TestLookupCurrent_InvalidTicker(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}
	err := lookupCurrent("ZZZZZZZZZZZ")
	if err == nil {
		t.Error("expected error for invalid ticker")
	}
}

// ---------------------------------------------------------------------------
// Date validation
// ---------------------------------------------------------------------------

func TestLookupDate_InvalidFormat(t *testing.T) {
	err := lookupDate("AAPL", "March 15 2024")
	if err == nil {
		t.Error("expected error for invalid date format")
	}
}

func TestLookupRange_InvalidFrom(t *testing.T) {
	err := lookupRange("AAPL", "bad-date", "", "1d")
	if err == nil {
		t.Error("expected error for invalid from date")
	}
}

func TestLookupRange_InvalidTo(t *testing.T) {
	err := lookupRange("AAPL", "2024-01-01", "bad-date", "1d")
	if err == nil {
		t.Error("expected error for invalid to date")
	}
}
