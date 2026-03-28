package convert

func init() {
	registerCategory(Category{
		Name: "Volume",
		Base: "ml",
		Units: map[string]Unit{
			"ml":          {Name: "milliliter", Factor: 1},
			"milliliter":  {Name: "milliliter", Factor: 1},
			"l":           {Name: "liter", Factor: 1000},
			"liter":       {Name: "liter", Factor: 1000},
			"gal":         {Name: "gallon (US)", Factor: 3785.41},
			"gallon":      {Name: "gallon (US)", Factor: 3785.41},
			"qt":          {Name: "quart (US)", Factor: 946.353},
			"quart":       {Name: "quart (US)", Factor: 946.353},
			"pt":          {Name: "pint (US)", Factor: 473.176},
			"pint":        {Name: "pint (US)", Factor: 473.176},
			"cup":         {Name: "cup (US)", Factor: 236.588},
			"floz":        {Name: "fluid ounce (US)", Factor: 29.5735},
			"tbsp":        {Name: "tablespoon", Factor: 14.7868},
			"tablespoon":  {Name: "tablespoon", Factor: 14.7868},
			"tsp":         {Name: "teaspoon", Factor: 4.92892},
			"teaspoon":    {Name: "teaspoon", Factor: 4.92892},
		},
	})
}
