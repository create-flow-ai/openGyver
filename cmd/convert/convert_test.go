package convert

import (
	"math"
	"testing"
)

func assertClose(t *testing.T, got, want float64, label string) {
	t.Helper()
	// Use relative tolerance for large values, absolute for small
	tol := 0.01
	if math.Abs(want) > 100 {
		tol = math.Abs(want) * 0.001 // 0.1% relative
	}
	if math.Abs(got-want) > tol {
		t.Errorf("%s: got %f, want %f", label, got, want)
	}
}

// ---------------------------------------------------------------------------
// Core helpers
// ---------------------------------------------------------------------------

func TestConvertByFactor(t *testing.T) {
	// 100 cm → in: 100 * 0.01 / 0.0254 = 39.3701
	from := Unit{Name: "centimeter", Factor: 0.01}
	to := Unit{Name: "inch", Factor: 0.0254}
	assertClose(t, convertByFactor(100, from, to), 39.3701, "100cm→in")
}

func TestConvertByFactor_Identity(t *testing.T) {
	u := Unit{Name: "meter", Factor: 1}
	assertClose(t, convertByFactor(42, u, u), 42, "identity conversion")
}

func TestLookup_Known(t *testing.T) {
	cat, ok := lookup("cm")
	if !ok {
		t.Fatal("expected to find 'cm'")
	}
	if cat.Name != "Length" {
		t.Errorf("expected Length category, got %s", cat.Name)
	}
}

func TestLookup_Unknown(t *testing.T) {
	_, ok := lookup("foobar")
	if ok {
		t.Error("expected lookup to fail for 'foobar'")
	}
}

func TestFormatNumber_Integer(t *testing.T) {
	if got := formatNumber(100); got != "100" {
		t.Errorf("got %q, want %q", got, "100")
	}
}

func TestFormatNumber_Decimal(t *testing.T) {
	got := formatNumber(39.370079)
	if got != "39.370079" {
		t.Errorf("got %q, want %q", got, "39.370079")
	}
}

func TestFormatNumber_TrailingZeros(t *testing.T) {
	got := formatNumber(1.5)
	if got != "1.5" {
		t.Errorf("got %q, want %q", got, "1.5")
	}
}

// ---------------------------------------------------------------------------
// Temperature
// ---------------------------------------------------------------------------

func TestTemperature_CtoF(t *testing.T) {
	cat, _ := lookup("c")
	result, err := cat.Convert(100, "c", "f")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 212, "100C→F")
}

func TestTemperature_FtoC(t *testing.T) {
	cat, _ := lookup("f")
	result, err := cat.Convert(72, "f", "c")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 22.22, "72F→C")
}

func TestTemperature_CtoK(t *testing.T) {
	cat, _ := lookup("c")
	result, err := cat.Convert(0, "c", "k")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 273.15, "0C→K")
}

func TestTemperature_KtoC(t *testing.T) {
	cat, _ := lookup("k")
	result, err := cat.Convert(300, "k", "c")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 26.85, "300K→C")
}

func TestTemperature_FtoK(t *testing.T) {
	cat, _ := lookup("f")
	result, err := cat.Convert(32, "f", "k")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 273.15, "32F→K")
}

func TestTemperature_KtoF(t *testing.T) {
	cat, _ := lookup("k")
	result, err := cat.Convert(273.15, "k", "f")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 32, "273.15K→F")
}

func TestTemperature_Aliases(t *testing.T) {
	cat, _ := lookup("celsius")
	result, err := cat.Convert(100, "celsius", "fahrenheit")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 212, "celsius alias")
}

func TestTemperature_Identity(t *testing.T) {
	cat, _ := lookup("c")
	result, err := cat.Convert(37, "c", "c")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, 37, "C→C identity")
}

func TestTemperature_AbsoluteZero(t *testing.T) {
	cat, _ := lookup("k")
	result, err := cat.Convert(0, "k", "c")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, -273.15, "0K→C")
}

func TestTemperature_NegativeF(t *testing.T) {
	cat, _ := lookup("f")
	result, err := cat.Convert(-40, "f", "c")
	if err != nil {
		t.Fatal(err)
	}
	assertClose(t, result, -40, "-40F→C")
}

// ---------------------------------------------------------------------------
// Length
// ---------------------------------------------------------------------------

func TestLength_CmToIn(t *testing.T) {
	cat, _ := lookup("cm")
	from, to := cat.Units["cm"], cat.Units["in"]
	assertClose(t, convertByFactor(100, from, to), 39.37, "100cm→in")
}

func TestLength_KmToMi(t *testing.T) {
	cat, _ := lookup("km")
	from, to := cat.Units["km"], cat.Units["mi"]
	assertClose(t, convertByFactor(5, from, to), 3.107, "5km→mi")
}

func TestLength_FtToM(t *testing.T) {
	cat, _ := lookup("ft")
	from, to := cat.Units["ft"], cat.Units["m"]
	assertClose(t, convertByFactor(6, from, to), 1.8288, "6ft→m")
}

func TestLength_MmToKm(t *testing.T) {
	cat, _ := lookup("mm")
	from, to := cat.Units["mm"], cat.Units["km"]
	assertClose(t, convertByFactor(1000000, from, to), 1, "1000000mm→km")
}

func TestLength_YdToFt(t *testing.T) {
	cat, _ := lookup("yd")
	from, to := cat.Units["yd"], cat.Units["ft"]
	assertClose(t, convertByFactor(1, from, to), 3, "1yd→ft")
}

func TestLength_NauticalMileToKm(t *testing.T) {
	cat, _ := lookup("nm")
	from, to := cat.Units["nm"], cat.Units["km"]
	assertClose(t, convertByFactor(1, from, to), 1.852, "1nm→km")
}

func TestLength_Aliases(t *testing.T) {
	cat, _ := lookup("feet")
	if cat.Name != "Length" {
		t.Errorf("expected Length, got %s", cat.Name)
	}
	cat2, _ := lookup("millimeter")
	if cat2.Name != "Length" {
		t.Errorf("expected Length for millimeter alias, got %s", cat2.Name)
	}
}

// ---------------------------------------------------------------------------
// Weight
// ---------------------------------------------------------------------------

func TestWeight_LbToKg(t *testing.T) {
	cat, _ := lookup("lb")
	from, to := cat.Units["lb"], cat.Units["kg"]
	assertClose(t, convertByFactor(150, from, to), 68.04, "150lb→kg")
}

func TestWeight_KgToOz(t *testing.T) {
	cat, _ := lookup("kg")
	from, to := cat.Units["kg"], cat.Units["oz"]
	assertClose(t, convertByFactor(1, from, to), 35.274, "1kg→oz")
}

func TestWeight_TonToTonne(t *testing.T) {
	cat, _ := lookup("ton")
	from, to := cat.Units["ton"], cat.Units["tonne"]
	assertClose(t, convertByFactor(1, from, to), 0.9072, "1ton→tonne")
}

func TestWeight_StToLb(t *testing.T) {
	cat, _ := lookup("st")
	from, to := cat.Units["st"], cat.Units["lb"]
	assertClose(t, convertByFactor(10, from, to), 140, "10st→lb")
}

func TestWeight_MgToG(t *testing.T) {
	cat, _ := lookup("mg")
	from, to := cat.Units["mg"], cat.Units["g"]
	assertClose(t, convertByFactor(1000, from, to), 1, "1000mg→g")
}

func TestWeight_Aliases(t *testing.T) {
	aliases := []string{"pound", "ounce", "kilogram", "milligram", "gram", "stone"}
	for _, a := range aliases {
		cat, ok := lookup(a)
		if !ok {
			t.Errorf("alias %q not found", a)
			continue
		}
		if cat.Name != "Weight" {
			t.Errorf("alias %q: expected Weight, got %s", a, cat.Name)
		}
	}
}

// ---------------------------------------------------------------------------
// Volume
// ---------------------------------------------------------------------------

func TestVolume_MlToCup(t *testing.T) {
	cat, _ := lookup("ml")
	from, to := cat.Units["ml"], cat.Units["cup"]
	assertClose(t, convertByFactor(500, from, to), 2.113, "500ml→cup")
}

func TestVolume_GalToL(t *testing.T) {
	cat, _ := lookup("gal")
	from, to := cat.Units["gal"], cat.Units["l"]
	assertClose(t, convertByFactor(1, from, to), 3.785, "1gal→l")
}

func TestVolume_TbspToTsp(t *testing.T) {
	cat, _ := lookup("tbsp")
	from, to := cat.Units["tbsp"], cat.Units["tsp"]
	assertClose(t, convertByFactor(1, from, to), 3, "1tbsp→tsp")
}

func TestVolume_CupToFloz(t *testing.T) {
	cat, _ := lookup("cup")
	from, to := cat.Units["cup"], cat.Units["floz"]
	assertClose(t, convertByFactor(1, from, to), 8, "1cup→floz")
}

func TestVolume_QtToPt(t *testing.T) {
	cat, _ := lookup("qt")
	from, to := cat.Units["qt"], cat.Units["pt"]
	assertClose(t, convertByFactor(1, from, to), 2, "1qt→pt")
}

func TestVolume_Aliases(t *testing.T) {
	aliases := []string{"milliliter", "liter", "gallon", "quart", "pint", "tablespoon", "teaspoon"}
	for _, a := range aliases {
		cat, ok := lookup(a)
		if !ok {
			t.Errorf("alias %q not found", a)
			continue
		}
		if cat.Name != "Volume" {
			t.Errorf("alias %q: expected Volume, got %s", a, cat.Name)
		}
	}
}

// ---------------------------------------------------------------------------
// Area
// ---------------------------------------------------------------------------

func TestArea_AcreToSqft(t *testing.T) {
	cat, _ := lookup("acre")
	from, to := cat.Units["acre"], cat.Units["sqft"]
	assertClose(t, convertByFactor(1, from, to), 43560, "1acre→sqft")
}

func TestArea_SqmiToAcre(t *testing.T) {
	cat, _ := lookup("sqmi")
	from, to := cat.Units["sqmi"], cat.Units["acre"]
	assertClose(t, convertByFactor(1, from, to), 640, "1sqmi→acre")
}

func TestArea_HectareToAcre(t *testing.T) {
	cat, _ := lookup("hectare")
	from, to := cat.Units["hectare"], cat.Units["acre"]
	assertClose(t, convertByFactor(1, from, to), 2.471, "1hectare→acre")
}

func TestArea_SqmToSqft(t *testing.T) {
	cat, _ := lookup("sqm")
	from, to := cat.Units["sqm"], cat.Units["sqft"]
	assertClose(t, convertByFactor(1, from, to), 10.764, "1sqm→sqft")
}

func TestArea_HaAlias(t *testing.T) {
	cat, ok := lookup("ha")
	if !ok {
		t.Fatal("ha alias not found")
	}
	if cat.Name != "Area" {
		t.Errorf("expected Area, got %s", cat.Name)
	}
}

// ---------------------------------------------------------------------------
// Speed
// ---------------------------------------------------------------------------

func TestSpeed_MphToKph(t *testing.T) {
	cat, _ := lookup("mph")
	from, to := cat.Units["mph"], cat.Units["kph"]
	assertClose(t, convertByFactor(60, from, to), 96.56, "60mph→kph")
}

func TestSpeed_KphToMph(t *testing.T) {
	cat, _ := lookup("kph")
	from, to := cat.Units["kph"], cat.Units["mph"]
	assertClose(t, convertByFactor(100, from, to), 62.14, "100kph→mph")
}

func TestSpeed_KnotsToMph(t *testing.T) {
	cat, _ := lookup("knots")
	from, to := cat.Units["knots"], cat.Units["mph"]
	assertClose(t, convertByFactor(30, from, to), 34.52, "30knots→mph")
}

func TestSpeed_MpsToKph(t *testing.T) {
	cat, _ := lookup("mps")
	from, to := cat.Units["mps"], cat.Units["kph"]
	assertClose(t, convertByFactor(10, from, to), 36, "10mps→kph")
}

func TestSpeed_Aliases(t *testing.T) {
	pairs := map[string]string{"m/s": "Speed", "km/h": "Speed", "ft/s": "Speed", "knot": "Speed"}
	for alias, want := range pairs {
		cat, ok := lookup(alias)
		if !ok {
			t.Errorf("alias %q not found", alias)
			continue
		}
		if cat.Name != want {
			t.Errorf("alias %q: expected %s, got %s", alias, want, cat.Name)
		}
	}
}

// ---------------------------------------------------------------------------
// Data
// ---------------------------------------------------------------------------

func TestData_GbToMb(t *testing.T) {
	cat, _ := lookup("gb")
	from, to := cat.Units["gb"], cat.Units["mb"]
	assertClose(t, convertByFactor(1.5, from, to), 1536, "1.5gb→mb")
}

func TestData_TbToGb(t *testing.T) {
	cat, _ := lookup("tb")
	from, to := cat.Units["tb"], cat.Units["gb"]
	assertClose(t, convertByFactor(1, from, to), 1024, "1tb→gb")
}

func TestData_MbitToMb(t *testing.T) {
	cat, _ := lookup("mbit")
	from, to := cat.Units["mbit"], cat.Units["mb"]
	assertClose(t, convertByFactor(100, from, to), 12.5, "100mbit→mb")
}

func TestData_BitToByte(t *testing.T) {
	cat, _ := lookup("bit")
	from, to := cat.Units["bit"], cat.Units["b"]
	assertClose(t, convertByFactor(8, from, to), 1, "8bit→b")
}

func TestData_KbToB(t *testing.T) {
	cat, _ := lookup("kb")
	from, to := cat.Units["kb"], cat.Units["b"]
	assertClose(t, convertByFactor(1, from, to), 1024, "1kb→b")
}

func TestData_PbToTb(t *testing.T) {
	cat, _ := lookup("pb")
	from, to := cat.Units["pb"], cat.Units["tb"]
	assertClose(t, convertByFactor(1, from, to), 1024, "1pb→tb")
}

// ---------------------------------------------------------------------------
// Time
// ---------------------------------------------------------------------------

func TestTime_DaysToHours(t *testing.T) {
	cat, _ := lookup("day")
	from, to := cat.Units["day"], cat.Units["hr"]
	assertClose(t, convertByFactor(365, from, to), 8760, "365days→hours")
}

func TestTime_YearToDays(t *testing.T) {
	cat, _ := lookup("year")
	from, to := cat.Units["year"], cat.Units["day"]
	assertClose(t, convertByFactor(1, from, to), 365, "1year→days")
}

func TestTime_MinToHr(t *testing.T) {
	cat, _ := lookup("min")
	from, to := cat.Units["min"], cat.Units["hr"]
	assertClose(t, convertByFactor(90, from, to), 1.5, "90min→hr")
}

func TestTime_MsToSec(t *testing.T) {
	cat, _ := lookup("ms")
	from, to := cat.Units["ms"], cat.Units["sec"]
	assertClose(t, convertByFactor(5000, from, to), 5, "5000ms→sec")
}

func TestTime_WeekToDay(t *testing.T) {
	cat, _ := lookup("week")
	from, to := cat.Units["week"], cat.Units["day"]
	assertClose(t, convertByFactor(2, from, to), 14, "2week→day")
}

func TestTime_MonthToWeek(t *testing.T) {
	cat, _ := lookup("month")
	from, to := cat.Units["month"], cat.Units["week"]
	assertClose(t, convertByFactor(1, from, to), 4.2857, "1month→week")
}

func TestTime_Aliases(t *testing.T) {
	aliases := []string{"millisecond", "second", "minute", "hour", "hours", "days", "weeks", "months", "years"}
	for _, a := range aliases {
		cat, ok := lookup(a)
		if !ok {
			t.Errorf("alias %q not found", a)
			continue
		}
		if cat.Name != "Time" {
			t.Errorf("alias %q: expected Time, got %s", a, cat.Name)
		}
	}
}

// ---------------------------------------------------------------------------
// Currency — unit registration only (no network calls)
// ---------------------------------------------------------------------------

func TestCurrency_Lookup(t *testing.T) {
	codes := []string{"usd", "eur", "gbp", "jpy", "cad", "aud", "chf", "cny", "inr", "mxn"}
	for _, c := range codes {
		cat, ok := lookup(c)
		if !ok {
			t.Errorf("currency %q not found", c)
			continue
		}
		if cat.Name != "Currency" {
			t.Errorf("currency %q: expected Currency, got %s", c, cat.Name)
		}
	}
}

func TestCurrency_AllRegistered(t *testing.T) {
	expected := len(currencyCodes)
	cat, _ := lookup("usd")
	if len(cat.Units) != expected {
		t.Errorf("expected %d currency units, got %d", expected, len(cat.Units))
	}
}

// ---------------------------------------------------------------------------
// Cross-category rejection
// ---------------------------------------------------------------------------

func TestCrossCategory_Rejected(t *testing.T) {
	cmCat, _ := lookup("cm")
	kgCat, _ := lookup("kg")
	if cmCat.Name == kgCat.Name {
		t.Error("cm and kg should be in different categories")
	}
}

// ---------------------------------------------------------------------------
// Registry completeness
// ---------------------------------------------------------------------------

func TestRegistry_AllCategoriesRegistered(t *testing.T) {
	expected := []string{"Temperature", "Length", "Weight", "Volume", "Area", "Speed", "Data", "Time", "Currency"}
	names := map[string]bool{}
	for _, cat := range registry {
		names[cat.Name] = true
	}
	for _, want := range expected {
		if !names[want] {
			t.Errorf("category %q not registered", want)
		}
	}
}
