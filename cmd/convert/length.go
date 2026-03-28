package convert

func init() {
	registerCategory(Category{
		Name: "Length",
		Base: "m",
		Units: map[string]Unit{
			"mm":         {Name: "millimeter", Factor: 0.001},
			"millimeter": {Name: "millimeter", Factor: 0.001},
			"cm":         {Name: "centimeter", Factor: 0.01},
			"centimeter": {Name: "centimeter", Factor: 0.01},
			"m":          {Name: "meter", Factor: 1},
			"meter":      {Name: "meter", Factor: 1},
			"km":         {Name: "kilometer", Factor: 1000},
			"kilometer":  {Name: "kilometer", Factor: 1000},
			"in":         {Name: "inch", Factor: 0.0254},
			"inch":       {Name: "inch", Factor: 0.0254},
			"ft":         {Name: "foot", Factor: 0.3048},
			"foot":       {Name: "foot", Factor: 0.3048},
			"feet":       {Name: "foot", Factor: 0.3048},
			"yd":         {Name: "yard", Factor: 0.9144},
			"yard":       {Name: "yard", Factor: 0.9144},
			"mi":         {Name: "mile", Factor: 1609.344},
			"mile":       {Name: "mile", Factor: 1609.344},
			"nm":         {Name: "nautical mile", Factor: 1852},
		},
	})
}
