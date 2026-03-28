package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	date     string
	from     string
	to       string
	lat      float64
	lon      float64
	units    string
	field    string
	jsonOut  bool
)

var weatherCmd = &cobra.Command{
	Use:   "weather <city>",
	Short: "Weather forecast, current conditions, and historical data",
	Long: `Look up weather for any city — current conditions, multi-day forecast,
or historical data back to 1940.

Uses Open-Meteo API (free, no API key required).

DATA RETURNED:

  Temperature (current, min, max), feels like, humidity, wind speed,
  wind direction, precipitation, cloud cover, pressure, UV index,
  sunrise/sunset, weather description.

OPTIONS:

  (no flags)         Current weather
  --date             Weather on a specific date (past or future up to 16 days)
  --from / --to      Date range (historical or forecast)
  --units            celsius (default) or fahrenheit
  --field / -f       Output a single value (temperature, humidity, wind, etc.)
  --lat / --lon      Use coordinates instead of city name

EXAMPLES:

  openGyver weather "New York"
  openGyver weather Tokyo --units fahrenheit
  openGyver weather London --date 2024-12-25
  openGyver weather Paris --from 2024-06-01 --to 2024-06-07
  openGyver weather "San Francisco" -f temperature
  openGyver weather --lat 40.7128 --lon -74.006
  openGyver weather Berlin -j`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWeather,
}

func runWeather(c *cobra.Command, args []string) error {
	var latitude, longitude float64
	var cityName string

	if lat != 0 || lon != 0 {
		latitude, longitude = lat, lon
		cityName = fmt.Sprintf("%.4f, %.4f", lat, lon)
	} else if len(args) > 0 {
		city, err := geocode(args[0])
		if err != nil {
			return err
		}
		latitude, longitude = city.Lat, city.Lon
		cityName = fmt.Sprintf("%s, %s", city.Name, city.Country)
	} else {
		return fmt.Errorf("provide a city name or use --lat/--lon")
	}

	tempUnit := "celsius"
	windUnit := "kmh"
	if strings.EqualFold(units, "fahrenheit") || strings.EqualFold(units, "f") {
		tempUnit = "fahrenheit"
		windUnit = "mph"
	}

	if date != "" {
		return lookupDate(cityName, latitude, longitude, date, tempUnit, windUnit)
	}
	if from != "" || to != "" {
		return lookupRange(cityName, latitude, longitude, from, to, tempUnit, windUnit)
	}
	return lookupCurrent(cityName, latitude, longitude, tempUnit, windUnit)
}

// --- Geocoding ---

type geoResult struct {
	Name    string  `json:"name"`
	Country string  `json:"country"`
	Lat     float64 `json:"latitude"`
	Lon     float64 `json:"longitude"`
	Admin1  string  `json:"admin1"`
}

type geoResponse struct {
	Results []geoResult `json:"results"`
}

func geocode(city string) (*geoResult, error) {
	u := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json",
		url.QueryEscape(city))
	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}
	defer resp.Body.Close()

	var result geoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("parsing geocoding response: %w", err)
	}
	if len(result.Results) == 0 {
		return nil, fmt.Errorf("city not found: %s", city)
	}
	return &result.Results[0], nil
}

// --- Current weather ---

type currentResponse struct {
	CurrentWeather struct {
		Temperature   float64 `json:"temperature"`
		Windspeed     float64 `json:"windspeed"`
		WindDirection float64 `json:"winddirection"`
		WeatherCode   int     `json:"weathercode"`
		IsDay         int     `json:"is_day"`
		Time          string  `json:"time"`
	} `json:"current_weather"`
	Hourly struct {
		RelativeHumidity []jsonFloat `json:"relative_humidity_2m"`
		ApparentTemp     []jsonFloat `json:"apparent_temperature"`
		Precipitation    []jsonFloat `json:"precipitation"`
		CloudCover       []jsonFloat `json:"cloud_cover"`
		Pressure         []jsonFloat `json:"surface_pressure"`
		UVIndex          []jsonFloat `json:"uv_index"`
	} `json:"hourly"`
	Daily struct {
		Sunrise []string    `json:"sunrise"`
		Sunset  []string    `json:"sunset"`
		TempMax []jsonFloat `json:"temperature_2m_max"`
		TempMin []jsonFloat `json:"temperature_2m_min"`
	} `json:"daily"`
}

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

func lookupCurrent(city string, lat, lon float64, tempUnit, windUnit string) error {
	now := time.Now()
	hourIdx := now.Hour()

	u := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
		"&current_weather=true"+
		"&hourly=relative_humidity_2m,apparent_temperature,precipitation,cloud_cover,surface_pressure,uv_index"+
		"&daily=sunrise,sunset,temperature_2m_max,temperature_2m_min"+
		"&temperature_unit=%s&wind_speed_unit=%s&timezone=auto&forecast_days=1",
		lat, lon, tempUnit, windUnit)

	resp, err := fetchJSON(u)
	if err != nil {
		return err
	}

	var data currentResponse
	if err := json.Unmarshal(resp, &data); err != nil {
		return fmt.Errorf("parsing weather: %w", err)
	}

	cw := data.CurrentWeather
	tempSuffix := "C"
	speedSuffix := "km/h"
	if tempUnit == "fahrenheit" {
		tempSuffix = "F"
		speedSuffix = "mph"
	}

	humidity := safeIndex(data.Hourly.RelativeHumidity, hourIdx)
	feelsLike := safeIndex(data.Hourly.ApparentTemp, hourIdx)
	precip := safeIndex(data.Hourly.Precipitation, hourIdx)
	clouds := safeIndex(data.Hourly.CloudCover, hourIdx)
	pressure := safeIndex(data.Hourly.Pressure, hourIdx)
	uv := safeIndex(data.Hourly.UVIndex, hourIdx)
	desc := weatherDescription(cw.WeatherCode)

	fields := map[string]interface{}{
		"city": city, "temperature": cw.Temperature, "feels_like": feelsLike,
		"humidity": humidity, "wind_speed": cw.Windspeed,
		"wind_direction": cw.WindDirection, "precipitation": precip,
		"cloud_cover": clouds, "pressure": pressure, "uv_index": uv,
		"description": desc, "time": cw.Time,
	}

	if len(data.Daily.TempMax) > 0 {
		fields["temp_max"] = float64(data.Daily.TempMax[0])
		fields["temp_min"] = float64(data.Daily.TempMin[0])
	}
	if len(data.Daily.Sunrise) > 0 {
		fields["sunrise"] = data.Daily.Sunrise[0]
		fields["sunset"] = data.Daily.Sunset[0]
	}

	if jsonOut {
		return cmd.PrintJSON(fields)
	}
	if field != "" {
		return printField(field, fields)
	}

	fmt.Printf("Location:      %s\n", city)
	fmt.Printf("Condition:     %s\n", desc)
	fmt.Printf("Temperature:   %.1f%s (feels like %.1f%s)\n", cw.Temperature, tempSuffix, feelsLike, tempSuffix)
	if len(data.Daily.TempMax) > 0 {
		fmt.Printf("High/Low:      %.1f%s / %.1f%s\n",
			float64(data.Daily.TempMax[0]), tempSuffix, float64(data.Daily.TempMin[0]), tempSuffix)
	}
	fmt.Printf("Humidity:      %.0f%%\n", humidity)
	fmt.Printf("Wind:          %.1f %s (%.0f)\n", cw.Windspeed, speedSuffix, cw.WindDirection)
	fmt.Printf("Precipitation: %.1f mm\n", precip)
	fmt.Printf("Cloud cover:   %.0f%%\n", clouds)
	fmt.Printf("Pressure:      %.0f hPa\n", pressure)
	fmt.Printf("UV Index:      %.1f\n", uv)
	if len(data.Daily.Sunrise) > 0 {
		fmt.Printf("Sunrise:       %s\n", formatTime(data.Daily.Sunrise[0]))
		fmt.Printf("Sunset:        %s\n", formatTime(data.Daily.Sunset[0]))
	}
	return nil
}

// --- Date lookup ---

type dailyResponse struct {
	Daily struct {
		Time          []string    `json:"time"`
		TempMax       []jsonFloat `json:"temperature_2m_max"`
		TempMin       []jsonFloat `json:"temperature_2m_min"`
		Precipitation []jsonFloat `json:"precipitation_sum"`
		WindMax       []jsonFloat `json:"wind_speed_10m_max"`
		WindDir       []jsonFloat `json:"wind_direction_10m_dominant"`
		Humidity      []jsonFloat `json:"relative_humidity_2m_mean"`  // historical only
		Sunrise       []string    `json:"sunrise"`
		Sunset        []string    `json:"sunset"`
		UVMax         []jsonFloat `json:"uv_index_max"`
		WeatherCode   []int       `json:"weather_code"`
	} `json:"daily"`
}

func lookupDate(city string, lat, lon float64, dateStr, tempUnit, windUnit string) error {
	return lookupRange(city, lat, lon, dateStr, dateStr, tempUnit, windUnit)
}

func lookupRange(city string, lat, lon float64, fromStr, toStr, tempUnit, windUnit string) error {
	now := time.Now()
	fromDate := now.AddDate(0, 0, -7)
	toDate := now.AddDate(0, 0, 7)

	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			return fmt.Errorf("invalid --from date (use YYYY-MM-DD): %w", err)
		}
		fromDate = t
	}
	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			return fmt.Errorf("invalid --to date (use YYYY-MM-DD): %w", err)
		}
		toDate = t
	}
	if fromStr != "" && toStr == "" {
		toDate = fromDate
	}

	// Decide: forecast or archive API
	dailyVars := "temperature_2m_max,temperature_2m_min,precipitation_sum,wind_speed_10m_max,wind_direction_10m_dominant,sunrise,sunset,uv_index_max,weather_code"

	var apiURL string
	if fromDate.Before(now.AddDate(0, 0, -5)) {
		// Historical
		apiURL = fmt.Sprintf("https://archive-api.open-meteo.com/v1/archive?latitude=%.4f&longitude=%.4f"+
			"&start_date=%s&end_date=%s&daily=%s&temperature_unit=%s&wind_speed_unit=%s&timezone=auto",
			lat, lon, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"), dailyVars, tempUnit, windUnit)
	} else {
		// Forecast
		apiURL = fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
			"&start_date=%s&end_date=%s&daily=%s&temperature_unit=%s&wind_speed_unit=%s&timezone=auto",
			lat, lon, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"), dailyVars, tempUnit, windUnit)
	}

	resp, err := fetchJSON(apiURL)
	if err != nil {
		return err
	}

	var data dailyResponse
	if err := json.Unmarshal(resp, &data); err != nil {
		return fmt.Errorf("parsing weather: %w", err)
	}

	if len(data.Daily.Time) == 0 {
		return fmt.Errorf("no weather data available for the requested dates")
	}

	type dayRow struct {
		Date          string  `json:"date"`
		TempMax       float64 `json:"temp_max"`
		TempMin       float64 `json:"temp_min"`
		Precipitation float64 `json:"precipitation_mm"`
		WindMax       float64 `json:"wind_max"`
		WindDir       float64 `json:"wind_direction"`
		UVMax         float64 `json:"uv_max"`
		Description   string  `json:"description"`
		Sunrise       string  `json:"sunrise"`
		Sunset        string  `json:"sunset"`
	}

	var rows []dayRow
	for i, d := range data.Daily.Time {
		code := 0
		if i < len(data.Daily.WeatherCode) {
			code = data.Daily.WeatherCode[i]
		}
		rows = append(rows, dayRow{
			Date:          d,
			TempMax:       safeIdxFloat(data.Daily.TempMax, i),
			TempMin:       safeIdxFloat(data.Daily.TempMin, i),
			Precipitation: safeIdxFloat(data.Daily.Precipitation, i),
			WindMax:       safeIdxFloat(data.Daily.WindMax, i),
			WindDir:       safeIdxFloat(data.Daily.WindDir, i),
			UVMax:         safeIdxFloat(data.Daily.UVMax, i),
			Description:   weatherDescription(code),
			Sunrise:       safeIdxStr(data.Daily.Sunrise, i),
			Sunset:        safeIdxStr(data.Daily.Sunset, i),
		})
	}

	if jsonOut {
		return cmd.PrintJSON(map[string]interface{}{
			"city": city, "data": rows,
		})
	}

	if field != "" {
		for _, r := range rows {
			fmap := map[string]interface{}{
				"temperature": r.TempMax, "temp_max": r.TempMax, "temp_min": r.TempMin,
				"precipitation": r.Precipitation, "wind": r.WindMax, "wind_speed": r.WindMax,
				"uv": r.UVMax, "description": r.Description, "date": r.Date,
			}
			val, ok := fmap[strings.ToLower(field)]
			if !ok {
				return fmt.Errorf("unknown field: %q", field)
			}
			switch v := val.(type) {
			case float64:
				fmt.Printf("%.1f\n", v)
			default:
				fmt.Println(v)
			}
		}
		return nil
	}

	tempSuffix := "C"
	speedSuffix := "km/h"
	if tempUnit == "fahrenheit" {
		tempSuffix = "F"
		speedSuffix = "mph"
	}

	fmt.Printf("Location: %s\n\n", city)

	if len(rows) == 1 {
		r := rows[0]
		fmt.Printf("Date:          %s\n", r.Date)
		fmt.Printf("Condition:     %s\n", r.Description)
		fmt.Printf("High/Low:      %.1f%s / %.1f%s\n", r.TempMax, tempSuffix, r.TempMin, tempSuffix)
		fmt.Printf("Precipitation: %.1f mm\n", r.Precipitation)
		fmt.Printf("Wind (max):    %.1f %s\n", r.WindMax, speedSuffix)
		fmt.Printf("UV Index:      %.1f\n", r.UVMax)
		if r.Sunrise != "" {
			fmt.Printf("Sunrise:       %s\n", formatTime(r.Sunrise))
			fmt.Printf("Sunset:        %s\n", formatTime(r.Sunset))
		}
	} else {
		fmt.Printf("%-12s %-20s %8s %8s %8s %8s %6s\n", "Date", "Condition", "High", "Low", "Precip", "Wind", "UV")
		fmt.Printf("%-12s %-20s %8s %8s %8s %8s %6s\n", "────────────", "────────────────────", "────────", "────────", "────────", "────────", "──────")
		for _, r := range rows {
			desc := r.Description
			if len(desc) > 20 {
				desc = desc[:17] + "..."
			}
			fmt.Printf("%-12s %-20s %7.1f%s %7.1f%s %6.1fmm %5.1f%s %5.1f\n",
				r.Date, desc, r.TempMax, tempSuffix, r.TempMin, tempSuffix,
				r.Precipitation, r.WindMax, speedSuffix, r.UVMax)
		}
	}
	return nil
}

// --- Helpers ---

func fetchJSON(u string) ([]byte, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("weather API request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}
	var buf []byte
	buf, err = readAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	var buf []byte
	tmp := make([]byte, 4096)
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			break
		}
	}
	return buf, nil
}

func safeIndex(arr []jsonFloat, idx int) float64 {
	if idx < len(arr) {
		return float64(arr[idx])
	}
	if len(arr) > 0 {
		return float64(arr[len(arr)-1])
	}
	return 0
}

func safeIdxFloat(arr []jsonFloat, idx int) float64 {
	if idx < len(arr) {
		return float64(arr[idx])
	}
	return 0
}

func safeIdxStr(arr []string, idx int) string {
	if idx < len(arr) {
		return arr[idx]
	}
	return ""
}

func formatTime(iso string) string {
	t, err := time.Parse("2006-01-02T15:04", iso)
	if err != nil {
		return iso
	}
	return t.Format("3:04 PM")
}

func printField(f string, data map[string]interface{}) error {
	val, ok := data[strings.ToLower(f)]
	if !ok {
		return fmt.Errorf("unknown field: %q\nAvailable: temperature, feels_like, humidity, wind_speed, wind_direction, precipitation, cloud_cover, pressure, uv_index, description, temp_max, temp_min, sunrise, sunset", f)
	}
	switch v := val.(type) {
	case float64:
		fmt.Printf("%.1f\n", v)
	default:
		fmt.Println(v)
	}
	return nil
}

// WMO Weather interpretation codes
func weatherDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly clear"
	case 2:
		return "Partly cloudy"
	case 3:
		return "Overcast"
	case 45, 48:
		return "Fog"
	case 51:
		return "Light drizzle"
	case 53:
		return "Moderate drizzle"
	case 55:
		return "Dense drizzle"
	case 56, 57:
		return "Freezing drizzle"
	case 61:
		return "Slight rain"
	case 63:
		return "Moderate rain"
	case 65:
		return "Heavy rain"
	case 66, 67:
		return "Freezing rain"
	case 71:
		return "Slight snow"
	case 73:
		return "Moderate snow"
	case 75:
		return "Heavy snow"
	case 77:
		return "Snow grains"
	case 80:
		return "Slight rain showers"
	case 81:
		return "Moderate rain showers"
	case 82:
		return "Violent rain showers"
	case 85:
		return "Slight snow showers"
	case 86:
		return "Heavy snow showers"
	case 95:
		return "Thunderstorm"
	case 96, 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown"
	}
}

func init() {
	weatherCmd.Flags().StringVarP(&date, "date", "d", "", "weather on a specific date (YYYY-MM-DD)")
	weatherCmd.Flags().StringVar(&from, "from", "", "start date for range (YYYY-MM-DD)")
	weatherCmd.Flags().StringVar(&to, "to", "", "end date for range (YYYY-MM-DD)")
	weatherCmd.Flags().Float64Var(&lat, "lat", 0, "latitude (use with --lon instead of city name)")
	weatherCmd.Flags().Float64Var(&lon, "lon", 0, "longitude (use with --lat instead of city name)")
	weatherCmd.Flags().StringVar(&units, "units", "celsius", "temperature units: celsius or fahrenheit")
	weatherCmd.Flags().StringVarP(&field, "field", "f", "", "output single field: temperature, humidity, wind_speed, etc.")
	weatherCmd.Flags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")
	cmd.Register(weatherCmd)
}
