package testdata

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRandInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		n := randInt(10)
		if n < 0 || n >= 10 {
			t.Errorf("randInt(10) = %d, expected 0-9", n)
		}
	}
}

func TestRandPick(t *testing.T) {
	list := []string{"a", "b", "c"}
	for i := 0; i < 50; i++ {
		picked := randPick(list)
		found := false
		for _, item := range list {
			if item == picked {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("randPick returned %q which is not in the list", picked)
		}
	}
}

func TestPersonGeneration(t *testing.T) {
	// Generate a person inline using the same logic as the command.
	first := randPick(firstNames)
	last := randPick(lastNames)
	name := first + " " + last
	email := strings.ToLower(first) + "." + strings.ToLower(last) + "@" + randPick(domains)
	city := randPick(cities)
	state := randPick(states)
	age := 18 + randInt(62)

	if name == "" || name == " " {
		t.Error("name should not be empty")
	}
	if !strings.Contains(email, "@") {
		t.Error("email should contain @")
	}
	if city == "" {
		t.Error("city should not be empty")
	}
	if state == "" {
		t.Error("state should not be empty")
	}
	if age < 18 || age > 79 {
		t.Errorf("age = %d, expected 18-79", age)
	}
}

func TestPersonGeneration_NonEmptyFields(t *testing.T) {
	for i := 0; i < 10; i++ {
		first := randPick(firstNames)
		last := randPick(lastNames)
		if first == "" {
			t.Error("first name is empty")
		}
		if last == "" {
			t.Error("last name is empty")
		}
		domain := randPick(domains)
		if !strings.Contains(domain, ".") {
			t.Errorf("domain %q doesn't contain a dot", domain)
		}
		street := randPick(streets)
		if street == "" {
			t.Error("street is empty")
		}
	}
}

func TestRandomNumberRange(t *testing.T) {
	min := 10
	max := 20
	for i := 0; i < 100; i++ {
		v := min + randInt(max-min+1)
		if v < min || v > max {
			t.Errorf("random number %d outside range [%d, %d]", v, min, max)
		}
	}
}

func TestGenColumn_Name(t *testing.T) {
	result := genColumn("name")
	if result == "" {
		t.Error("genColumn(name) returned empty string")
	}
	parts := strings.Fields(result)
	if len(parts) != 2 {
		t.Errorf("genColumn(name) = %q, expected first+last name", result)
	}
}

func TestGenColumn_Email(t *testing.T) {
	result := genColumn("email")
	if !strings.Contains(result, "@") {
		t.Errorf("genColumn(email) = %q, expected to contain @", result)
	}
}

func TestGenColumn_Number(t *testing.T) {
	result := genColumn("number")
	n, err := strconv.Atoi(result)
	if err != nil {
		t.Errorf("genColumn(number) = %q, expected a number", result)
	}
	if n < 18 || n > 79 {
		t.Errorf("genColumn(number) = %d, expected 18-79", n)
	}
}

func TestGenColumn_Age(t *testing.T) {
	result := genColumn("age")
	n, err := strconv.Atoi(result)
	if err != nil {
		t.Errorf("genColumn(age) = %q, expected a number", result)
	}
	if n < 18 || n > 79 {
		t.Errorf("genColumn(age) = %d, expected 18-79", n)
	}
}

func TestGenColumn_Date(t *testing.T) {
	result := genColumn("date")
	_, err := time.Parse("2006-01-02", result)
	if err != nil {
		t.Errorf("genColumn(date) = %q, not a valid date: %v", result, err)
	}
}

func TestGenColumn_Bool(t *testing.T) {
	result := genColumn("bool")
	if result != "true" && result != "false" {
		t.Errorf("genColumn(bool) = %q, expected true or false", result)
	}
}

func TestGenColumn_City(t *testing.T) {
	result := genColumn("city")
	if result == "" {
		t.Error("genColumn(city) returned empty string")
	}
}

func TestGenColumn_Country(t *testing.T) {
	result := genColumn("country")
	if len(result) != 2 {
		t.Errorf("genColumn(country) = %q, expected 2-letter code", result)
	}
}

func TestGenColumn_Phone(t *testing.T) {
	result := genColumn("phone")
	if !strings.Contains(result, "(") || !strings.Contains(result, ")") {
		t.Errorf("genColumn(phone) = %q, expected phone format", result)
	}
}

func TestGenColumn_UUID(t *testing.T) {
	result := genColumn("uuid")
	parts := strings.Split(result, "-")
	if len(parts) != 5 {
		t.Errorf("genColumn(uuid) = %q, expected 5 dash-separated parts", result)
	}
}

func TestGenColumn_Unknown(t *testing.T) {
	result := genColumn("unknown_type")
	if !strings.HasPrefix(result, "val_") {
		t.Errorf("genColumn(unknown) = %q, expected val_ prefix", result)
	}
}
