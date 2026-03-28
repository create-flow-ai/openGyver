package convert

import "fmt"

func init() {
	registerCategory(Category{
		Name: "Temperature",
		Base: "c",
		Units: map[string]Unit{
			"c":          {Name: "Celsius"},
			"celsius":    {Name: "Celsius"},
			"f":          {Name: "Fahrenheit"},
			"fahrenheit": {Name: "Fahrenheit"},
			"k":          {Name: "Kelvin"},
			"kelvin":     {Name: "Kelvin"},
		},
		Convert: func(val float64, from, to string) (float64, error) {
			// Normalize aliases to base names
			norm := map[string]string{
				"c": "c", "celsius": "c",
				"f": "f", "fahrenheit": "f",
				"k": "k", "kelvin": "k",
			}
			f, t := norm[from], norm[to]

			// Convert to Celsius first
			var inC float64
			switch f {
			case "c":
				inC = val
			case "f":
				inC = (val - 32) * 5 / 9
			case "k":
				inC = val - 273.15
			default:
				return 0, fmt.Errorf("unknown temperature unit: %s", from)
			}

			// Convert from Celsius to target
			switch t {
			case "c":
				return inC, nil
			case "f":
				return inC*9/5 + 32, nil
			case "k":
				return inC + 273.15, nil
			default:
				return 0, fmt.Errorf("unknown temperature unit: %s", to)
			}
		},
	})
}
