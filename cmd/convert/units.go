package convert

import "fmt"

// Unit defines a convertible unit with its canonical name and factor to the category's base unit.
type Unit struct {
	Name   string
	Factor float64 // 1 of this unit = Factor base units (unused for temperature)
}

// Category groups related units and knows how to convert between them.
type Category struct {
	Name    string
	Base    string            // base unit name (e.g. "m" for length)
	Units   map[string]Unit   // alias → Unit
	Convert func(val float64, from, to string) (float64, error) // custom converter (for temperature)
}

// registry holds all registered categories.
var registry []Category

func registerCategory(c Category) {
	registry = append(registry, c)
}

// lookup finds which category a unit alias belongs to and returns it.
func lookup(alias string) (*Category, bool) {
	for i := range registry {
		if _, ok := registry[i].Units[alias]; ok {
			return &registry[i], true
		}
	}
	return nil, false
}

// convert performs a unit conversion using the standard factor-based approach:
//
//	value * fromFactor / toFactor
func convertByFactor(val float64, from, to Unit) float64 {
	return val * from.Factor / to.Factor
}

// listUnits returns a formatted string of all units in a category.
func listUnits(c *Category) string {
	var s string
	seen := map[string]bool{}
	for alias, u := range c.Units {
		if !seen[u.Name] {
			s += fmt.Sprintf("    %-8s %s\n", alias, u.Name)
			seen[u.Name] = true
		}
	}
	return s
}
