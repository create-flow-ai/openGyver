package convert

func init() {
	registerCategory(Category{
		Name: "Area",
		Base: "sqm",
		Units: map[string]Unit{
			"sqmm":    {Name: "square millimeter", Factor: 0.000001},
			"sqcm":    {Name: "square centimeter", Factor: 0.0001},
			"sqm":     {Name: "square meter", Factor: 1},
			"sqkm":    {Name: "square kilometer", Factor: 1000000},
			"sqin":    {Name: "square inch", Factor: 0.00064516},
			"sqft":    {Name: "square foot", Factor: 0.092903},
			"sqyd":    {Name: "square yard", Factor: 0.836127},
			"sqmi":    {Name: "square mile", Factor: 2589988},
			"acre":    {Name: "acre", Factor: 4046.86},
			"hectare": {Name: "hectare", Factor: 10000},
			"ha":      {Name: "hectare", Factor: 10000},
		},
	})
}
