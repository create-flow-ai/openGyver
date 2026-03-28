package convert

func init() {
	registerCategory(Category{
		Name: "Time",
		Base: "sec",
		Units: map[string]Unit{
			"ms":           {Name: "millisecond", Factor: 0.001},
			"millisecond":  {Name: "millisecond", Factor: 0.001},
			"sec":          {Name: "second", Factor: 1},
			"second":       {Name: "second", Factor: 1},
			"min":          {Name: "minute", Factor: 60},
			"minute":       {Name: "minute", Factor: 60},
			"hr":           {Name: "hour", Factor: 3600},
			"hour":         {Name: "hour", Factor: 3600},
			"hours":        {Name: "hour", Factor: 3600},
			"day":          {Name: "day", Factor: 86400},
			"days":         {Name: "day", Factor: 86400},
			"week":         {Name: "week", Factor: 604800},
			"weeks":        {Name: "week", Factor: 604800},
			"month":        {Name: "month (30d)", Factor: 2592000},
			"months":       {Name: "month (30d)", Factor: 2592000},
			"year":         {Name: "year (365d)", Factor: 31536000},
			"years":        {Name: "year (365d)", Factor: 31536000},
		},
	})
}
