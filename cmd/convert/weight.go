package convert

func init() {
	registerCategory(Category{
		Name: "Weight",
		Base: "g",
		Units: map[string]Unit{
			"mg":        {Name: "milligram", Factor: 0.001},
			"milligram": {Name: "milligram", Factor: 0.001},
			"g":         {Name: "gram", Factor: 1},
			"gram":      {Name: "gram", Factor: 1},
			"kg":        {Name: "kilogram", Factor: 1000},
			"kilogram":  {Name: "kilogram", Factor: 1000},
			"oz":        {Name: "ounce", Factor: 28.3495},
			"ounce":     {Name: "ounce", Factor: 28.3495},
			"lb":        {Name: "pound", Factor: 453.592},
			"pound":     {Name: "pound", Factor: 453.592},
			"ton":       {Name: "short ton", Factor: 907185},
			"tonne":     {Name: "metric tonne", Factor: 1000000},
			"st":        {Name: "stone", Factor: 6350.29},
			"stone":     {Name: "stone", Factor: 6350.29},
		},
	})
}
