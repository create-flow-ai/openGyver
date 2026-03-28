package stock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	date     string
	from     string
	to       string
	market   string
	interval string
	field    string
	jsonOut  bool
)

// Market suffix mappings for convenience.
var marketSuffixes = map[string]string{
	"nasdaq":    "",
	"nyse":      "",
	"us":        "",
	"kosdaq":    ".KQ",
	"kospi":     ".KS",
	"korea":     ".KS",
	"tokyo":     ".T",
	"tse":       ".T",
	"japan":     ".T",
	"london":    ".L",
	"lse":       ".L",
	"uk":        ".L",
	"hongkong":  ".HK",
	"hkex":      ".HK",
	"hk":        ".HK",
	"shanghai":  ".SS",
	"shenzhen":  ".SZ",
	"frankfurt":  ".DE",
	"xetra":     ".DE",
	"germany":   ".DE",
	"paris":     ".PA",
	"euronext":  ".PA",
	"france":    ".PA",
	"toronto":   ".TO",
	"tsx":       ".TO",
	"canada":    ".TO",
	"australia": ".AX",
	"asx":       ".AX",
	"mumbai":    ".BO",
	"bse":       ".BO",
	"nse":       ".NS",
	"india":     ".NS",
	"taiwan":    ".TW",
	"twse":      ".TW",
	"singapore": ".SI",
	"sgx":       ".SI",
	"brazil":    ".SA",
	"bovespa":   ".SA",
	"mexico":    ".MX",
	"swiss":     ".SW",
	"six":       ".SW",
	"amsterdam": ".AS",
	"stockholm": ".ST",
	"oslo":      ".OL",
	"copenhagen": ".CO",
	"helsinki":   ".HE",
	"jakarta":   ".JK",
	"bangkok":   ".BK",
	"johannesburg": ".JO",
	"newzealand": ".NZ",
}

var stockCmd = &cobra.Command{
	Use:   "stock <ticker>",
	Short: "Look up stock prices by ticker symbol",
	Long: `Look up current or historical stock prices from global markets.

Uses Yahoo Finance data (no API key required). Tickers are searched
universally — just type the symbol and it auto-resolves. Use --market
to target a specific exchange.

TICKER FORMAT:

  US stocks:     AAPL, MSFT, GOOGL, TSLA
  With suffix:   005930.KS (Samsung on KOSPI), 7203.T (Toyota on TSE)
  Or use --market: openGyver stock 005930 --market kospi

MARKET NAMES (for --market flag):

  US:      nasdaq, nyse, us
  Korea:   kosdaq, kospi, korea
  Japan:   tokyo, tse, japan
  UK:      london, lse, uk
  China:   shanghai, shenzhen, hongkong, hk
  Europe:  frankfurt, xetra, paris, euronext, amsterdam, swiss
  Americas: toronto, tsx, brazil, bovespa, mexico
  Asia:    singapore, sgx, taiwan, mumbai, nse, jakarta, bangkok
  Other:   australia, asx, johannesburg, oslo, stockholm, copenhagen,
           helsinki, newzealand

DATE OPTIONS:

  --date         Price on a specific date (YYYY-MM-DD)
  --from / --to  Historical range (defaults: --from 30 days ago, --to today)
  --interval     Data granularity: 1d (default), 1wk, 1mo

ABBREVIATED OUTPUT (--field / -f):

  Returns only the requested value — ideal for piping into other tools.

  price          Current / closing price
  change         Price change (absolute)
  percent        Price change (percentage)
  open           Opening price
  high           High price
  low            Low price
  close          Closing price
  volume         Trading volume
  currency       Currency code
  exchange       Exchange name

  Example: openGyver stock AAPL -f price   →   248.80

Examples:
  openGyver stock AAPL
  openGyver stock MSFT --date 2024-01-15
  openGyver stock AAPL --from 2024-01-01 --to 2024-06-30
  openGyver stock AAPL --from 2024-01-01 --to 2024-06-30 --interval 1wk
  openGyver stock 005930 --market kospi
  openGyver stock 7203 --market tokyo
  openGyver stock SHEL --market london
  openGyver stock 0700 --market hk
  openGyver stock AAPL -f price
  openGyver stock AAPL -f percent
  openGyver stock AAPL -f change`,
	Args: cobra.ExactArgs(1),
	RunE: runStock,
}

func runStock(c *cobra.Command, args []string) error {
	ticker := resolveTicker(args[0], market)

	if date != "" {
		return lookupDate(ticker, date)
	}
	if from != "" || to != "" {
		return lookupRange(ticker, from, to, interval)
	}
	return lookupCurrent(ticker)
}

func resolveTicker(ticker, mkt string) string {
	ticker = strings.ToUpper(ticker)
	if mkt != "" {
		suffix, ok := marketSuffixes[strings.ToLower(mkt)]
		if ok && !strings.Contains(ticker, ".") {
			ticker += suffix
		}
	}
	return ticker
}

// --- Yahoo Finance API ---

const yahooBaseURL = "https://query1.finance.yahoo.com/v8/finance/chart"

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol             string  `json:"symbol"`
				Currency           string  `json:"currency"`
				ExchangeName       string  `json:"exchangeName"`
				FullExchangeName   string  `json:"fullExchangeName"`
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				PreviousClose      float64 `json:"previousClose"`
				RegularMarketTime  int64   `json:"regularMarketTime"`
				Timezone           string  `json:"exchangeTimezoneName"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []jsonFloat `json:"open"`
					High   []jsonFloat `json:"high"`
					Low    []jsonFloat `json:"low"`
					Close  []jsonFloat `json:"close"`
					Volume []jsonFloat `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

// jsonFloat handles null JSON numbers.
type jsonFloat float64

func (f *jsonFloat) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*f = 0
		return nil
	}
	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*f = jsonFloat(v)
	return nil
}

func fetchYahoo(url string) (*yahooResponse, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "openGyver/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Yahoo Finance returned status %d for the requested ticker", resp.StatusCode)
	}

	var result yahooResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	if result.Chart.Error != nil {
		return nil, fmt.Errorf("Yahoo Finance error: %s", result.Chart.Error.Description)
	}
	if len(result.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data found for ticker")
	}

	return &result, nil
}

// --- Lookup modes ---

func lookupCurrent(ticker string) error {
	url := fmt.Sprintf("%s/%s?range=1d&interval=5m", yahooBaseURL, ticker)
	data, err := fetchYahoo(url)
	if err != nil {
		return err
	}

	meta := data.Chart.Result[0].Meta
	price := meta.RegularMarketPrice
	prevClose := meta.PreviousClose
	change := price - prevClose
	changePct := 0.0
	if prevClose > 0 {
		changePct = (change / prevClose) * 100
	}

	t := time.Unix(meta.RegularMarketTime, 0)
	loc, _ := time.LoadLocation(meta.Timezone)
	if loc != nil {
		t = t.In(loc)
	}

	fields := map[string]interface{}{
		"symbol": meta.Symbol, "exchange": meta.FullExchangeName,
		"currency": meta.Currency, "price": price, "change": change,
		"percent": changePct, "previous_close": prevClose,
		"as_of": t.Format("2006-01-02T15:04:05Z07:00"),
	}

	if jsonOut {
		return cmd.PrintJSON(fields)
	}

	if field != "" {
		return printField(field, fields)
	}

	arrow := "▲"
	if change < 0 {
		arrow = "▼"
	}

	fmt.Printf("Symbol:    %s\n", meta.Symbol)
	fmt.Printf("Exchange:  %s\n", meta.FullExchangeName)
	fmt.Printf("Currency:  %s\n", meta.Currency)
	fmt.Printf("Price:     %.2f\n", price)
	fmt.Printf("Change:    %s %.2f (%.2f%%)\n", arrow, change, changePct)
	fmt.Printf("Prev Close: %.2f\n", prevClose)
	fmt.Printf("As of:     %s\n", t.Format("2006-01-02 15:04 MST"))

	return nil
}

func lookupDate(ticker, dateStr string) error {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format (use YYYY-MM-DD): %w", err)
	}

	p1 := t.Unix()
	p2 := t.Add(24 * time.Hour).Unix()

	url := fmt.Sprintf("%s/%s?period1=%d&period2=%d&interval=1d", yahooBaseURL, ticker, p1, p2)
	data, err := fetchYahoo(url)
	if err != nil {
		return err
	}

	result := data.Chart.Result[0]
	meta := result.Meta

	if len(result.Indicators.Quote) == 0 || len(result.Indicators.Quote[0].Close) == 0 {
		return fmt.Errorf("no trading data for %s on %s (market may have been closed)", ticker, dateStr)
	}

	q := result.Indicators.Quote[0]

	dateFields := map[string]interface{}{
		"symbol": meta.Symbol, "exchange": meta.FullExchangeName,
		"currency": meta.Currency, "date": dateStr,
		"open": float64(q.Open[0]), "high": float64(q.High[0]),
		"low": float64(q.Low[0]), "close": float64(q.Close[0]),
		"price": float64(q.Close[0]), "volume": float64(q.Volume[0]),
		"change": 0.0, "percent": 0.0,
	}

	if jsonOut {
		return cmd.PrintJSON(dateFields)
	}

	if field != "" {
		return printField(field, dateFields)
	}

	fmt.Printf("Symbol:    %s\n", meta.Symbol)
	fmt.Printf("Exchange:  %s\n", meta.FullExchangeName)
	fmt.Printf("Currency:  %s\n", meta.Currency)
	fmt.Printf("Date:      %s\n", dateStr)
	fmt.Printf("Open:      %.2f\n", float64(q.Open[0]))
	fmt.Printf("High:      %.2f\n", float64(q.High[0]))
	fmt.Printf("Low:       %.2f\n", float64(q.Low[0]))
	fmt.Printf("Close:     %.2f\n", float64(q.Close[0]))
	fmt.Printf("Volume:    %.0f\n", float64(q.Volume[0]))

	return nil
}

func lookupRange(ticker, fromStr, toStr, ivl string) error {
	now := time.Now()
	fromTime := now.AddDate(0, 0, -30)
	toTime := now

	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			return fmt.Errorf("invalid --from date (use YYYY-MM-DD): %w", err)
		}
		fromTime = t
	}
	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			return fmt.Errorf("invalid --to date (use YYYY-MM-DD): %w", err)
		}
		toTime = t.Add(24 * time.Hour)
	}
	if ivl == "" {
		ivl = "1d"
	}

	url := fmt.Sprintf("%s/%s?period1=%d&period2=%d&interval=%s",
		yahooBaseURL, ticker, fromTime.Unix(), toTime.Unix(), ivl)
	data, err := fetchYahoo(url)
	if err != nil {
		return err
	}

	result := data.Chart.Result[0]
	meta := result.Meta

	if len(result.Indicators.Quote) == 0 {
		return fmt.Errorf("no trading data for the requested period")
	}

	q := result.Indicators.Quote[0]

	// Build rows for JSON and field output
	type row struct {
		Date   string  `json:"date"`
		Open   float64 `json:"open"`
		High   float64 `json:"high"`
		Low    float64 `json:"low"`
		Close  float64 `json:"close"`
		Volume float64 `json:"volume"`
	}
	var rows []row
	for i, ts := range result.Timestamp {
		if i >= len(q.Close) {
			break
		}
		rows = append(rows, row{
			Date:   time.Unix(ts, 0).Format("2006-01-02"),
			Open:   float64(q.Open[i]),
			High:   float64(q.High[i]),
			Low:    float64(q.Low[i]),
			Close:  float64(q.Close[i]),
			Volume: float64(q.Volume[i]),
		})
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"symbol":   meta.Symbol,
			"exchange": meta.FullExchangeName,
			"currency": meta.Currency,
			"interval": ivl,
			"data":     rows,
		})
	}

	// Abbreviated: output just the requested field values, one per line
	if field != "" {
		for _, r := range rows {
			fmap := map[string]interface{}{
				"price": r.Close, "close": r.Close, "open": r.Open,
				"high": r.High, "low": r.Low, "volume": r.Volume,
				"date": r.Date,
			}
			val, ok := fmap[strings.ToLower(field)]
			if !ok {
				return fmt.Errorf("unknown field: %q", field)
			}
			switch v := val.(type) {
			case float64:
				fmt.Printf("%.2f\n", v)
			default:
				fmt.Println(v)
			}
		}
		return nil
	}

	fmt.Printf("Symbol:    %s\n", meta.Symbol)
	fmt.Printf("Exchange:  %s\n", meta.FullExchangeName)
	fmt.Printf("Currency:  %s\n", meta.Currency)
	fmt.Printf("Interval:  %s\n", ivl)
	fmt.Println()
	fmt.Printf("%-12s %10s %10s %10s %10s %15s\n", "Date", "Open", "High", "Low", "Close", "Volume")
	fmt.Printf("%-12s %10s %10s %10s %10s %15s\n", "────────────", "──────────", "──────────", "──────────", "──────────", "───────────────")

	for _, r := range rows {
		fmt.Printf("%-12s %10.2f %10.2f %10.2f %10.2f %15.0f\n",
			r.Date, r.Open, r.High, r.Low, r.Close, r.Volume)
	}

	return nil
}

// printField outputs a single field value for piping.
func printField(f string, data map[string]interface{}) error {
	f = strings.ToLower(f)
	val, ok := data[f]
	if !ok {
		return fmt.Errorf("unknown field: %q\nAvailable: price, change, percent, open, high, low, close, volume, currency, exchange", f)
	}
	switch v := val.(type) {
	case float64:
		if v == float64(int64(v)) && v > 1000 {
			fmt.Printf("%.0f\n", v)
		} else {
			fmt.Printf("%.2f\n", v)
		}
	case string:
		fmt.Println(v)
	default:
		fmt.Println(v)
	}
	return nil
}

func init() {
	stockCmd.Flags().StringVarP(&date, "date", "d", "", "look up price on a specific date (YYYY-MM-DD)")
	stockCmd.Flags().StringVar(&from, "from", "", "start date for historical range (YYYY-MM-DD)")
	stockCmd.Flags().StringVar(&to, "to", "", "end date for historical range (YYYY-MM-DD)")
	stockCmd.Flags().StringVarP(&market, "market", "m", "", "target market/exchange (e.g. kosdaq, tokyo, london)")
	stockCmd.Flags().StringVar(&interval, "interval", "", "data interval: 1d, 1wk, 1mo (default: 1d)")
	stockCmd.Flags().StringVarP(&field, "field", "f", "", "output a single field: price, change, percent, open, high, low, close, volume, currency, exchange")
	stockCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	cmd.Register(stockCmd)
}
