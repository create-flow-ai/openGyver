package convert

func init() {
	registerCategory(Category{
		Name: "Speed",
		Base: "m/s",
		Units: map[string]Unit{
			"mps":   {Name: "meters/sec", Factor: 1},
			"m/s":   {Name: "meters/sec", Factor: 1},
			"kph":   {Name: "km/hour", Factor: 0.277778},
			"km/h":  {Name: "km/hour", Factor: 0.277778},
			"mph":   {Name: "miles/hour", Factor: 0.44704},
			"knot":  {Name: "knot", Factor: 0.514444},
			"knots": {Name: "knot", Factor: 0.514444},
			"fps":   {Name: "feet/sec", Factor: 0.3048},
			"ft/s":  {Name: "feet/sec", Factor: 0.3048},
		},
	})
}
