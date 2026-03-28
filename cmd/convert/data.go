package convert

func init() {
	registerCategory(Category{
		Name: "Data",
		Base: "b",
		Units: map[string]Unit{
			"b":  {Name: "byte", Factor: 1},
			"kb": {Name: "kilobyte", Factor: 1024},
			"mb": {Name: "megabyte", Factor: 1048576},
			"gb": {Name: "gigabyte", Factor: 1073741824},
			"tb": {Name: "terabyte", Factor: 1099511627776},
			"pb": {Name: "petabyte", Factor: 1125899906842624},
			"bit":  {Name: "bit", Factor: 0.125},
			"kbit": {Name: "kilobit", Factor: 128},
			"mbit": {Name: "megabit", Factor: 131072},
			"gbit": {Name: "gigabit", Factor: 134217728},
		},
	})
}
