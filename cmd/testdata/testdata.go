package testdata

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var jsonOut bool

var testdataCmd = &cobra.Command{
	Use:   "testdata",
	Short: "Generate fake/test data — people, CSV, JSON, numbers",
	Long: `Generate random test data for development and testing.

SUBCOMMANDS:
  person    Generate fake person records (name, email, phone, address)
  csv       Generate sample CSV with configurable column types
  json      Generate sample JSON objects
  number    Generate random numbers

All subcommands support --json/-j for structured JSON output.

Examples:
  openGyver testdata person --count 5
  openGyver testdata csv --rows 20 --columns name,email,age,city
  openGyver testdata json --count 10
  openGyver testdata number --min 1 --max 100 --count 5`,
}

// --- Data pools ---
var firstNames = []string{
	"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda",
	"William", "Elizabeth", "David", "Barbara", "Richard", "Susan", "Joseph", "Jessica",
	"Thomas", "Sarah", "Charles", "Karen", "Daniel", "Lisa", "Matthew", "Nancy",
	"Anthony", "Betty", "Mark", "Margaret", "Donald", "Sandra", "Steven", "Ashley",
	"Paul", "Emily", "Andrew", "Donna", "Joshua", "Michelle", "Kenneth", "Carol",
	"Kevin", "Amanda", "Brian", "Dorothy", "George", "Melissa", "Timothy", "Deborah",
}
var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis",
	"Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson",
	"Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee", "Perez", "Thompson",
	"White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker",
	"Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill",
}
var streets = []string{
	"Main St", "Oak Ave", "Maple Dr", "Cedar Ln", "Pine St", "Elm St", "Washington Ave",
	"Park Rd", "Lake Dr", "Hill St", "River Rd", "Forest Ave", "Sunset Blvd", "Spring St",
	"Valley Rd", "Highland Ave", "Meadow Ln", "Church St", "Center St", "Broadway",
}
var cities = []string{
	"New York", "Los Angeles", "Chicago", "Houston", "Phoenix", "Philadelphia",
	"San Antonio", "San Diego", "Dallas", "Austin", "Portland", "Denver",
	"Seattle", "Boston", "Atlanta", "Miami", "Nashville", "Minneapolis",
}
var states = []string{
	"NY", "CA", "IL", "TX", "AZ", "PA", "FL", "OH", "GA", "NC",
	"WA", "CO", "OR", "MA", "TN", "MN", "VA", "MD", "IN", "WI",
}
var domains = []string{
	"gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "example.com",
	"mail.com", "proton.me", "icloud.com", "aol.com", "zoho.com",
}

func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func randPick(list []string) string { return list[randInt(len(list))] }

// --- Person subcommand ---
var personCount int

var personCmd = &cobra.Command{
	Use:   "person",
	Short: "Generate fake person records",
	Long: `Generate random person data with name, email, phone, and address.

Examples:
  openGyver testdata person
  openGyver testdata person --count 5
  openGyver testdata person --count 3 -j`,
	RunE: func(c *cobra.Command, args []string) error {
		type Person struct {
			Name    string `json:"name"`
			Email   string `json:"email"`
			Phone   string `json:"phone"`
			Address string `json:"address"`
			City    string `json:"city"`
			State   string `json:"state"`
			Zip     string `json:"zip"`
			Age     int    `json:"age"`
		}
		var people []Person
		for i := 0; i < personCount; i++ {
			first := randPick(firstNames)
			last := randPick(lastNames)
			p := Person{
				Name:    first + " " + last,
				Email:   fmt.Sprintf("%s.%s@%s", strings.ToLower(first), strings.ToLower(last), randPick(domains)),
				Phone:   fmt.Sprintf("(%03d) %03d-%04d", 200+randInt(800), 200+randInt(800), randInt(10000)),
				Address: fmt.Sprintf("%d %s", 100+randInt(9900), randPick(streets)),
				City:    randPick(cities),
				State:   randPick(states),
				Zip:     fmt.Sprintf("%05d", 10000+randInt(90000)),
				Age:     18 + randInt(62),
			}
			people = append(people, p)
		}
		if jsonOut {
			return cmd.PrintJSON(people)
		}
		for _, p := range people {
			fmt.Printf("%s | %s | %s | %s, %s, %s %s | Age: %d\n",
				p.Name, p.Email, p.Phone, p.Address, p.City, p.State, p.Zip, p.Age)
		}
		return nil
	},
}

// --- CSV subcommand ---
var (
	csvRows    int
	csvColumns string
)

var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Generate sample CSV data",
	Long: `Generate CSV with configurable columns.

Column types: name, email, number, date, bool, city, country, age, phone, uuid

Examples:
  openGyver testdata csv --rows 10 --columns name,email,age
  openGyver testdata csv --rows 5 --columns name,city,phone,date`,
	RunE: func(c *cobra.Command, args []string) error {
		cols := strings.Split(csvColumns, ",")
		w := csv.NewWriter(os.Stdout)
		w.Write(cols)
		for i := 0; i < csvRows; i++ {
			var row []string
			for _, col := range cols {
				row = append(row, genColumn(strings.TrimSpace(col)))
			}
			w.Write(row)
		}
		w.Flush()
		return nil
	},
}

func genColumn(colType string) string {
	switch strings.ToLower(colType) {
	case "name":
		return randPick(firstNames) + " " + randPick(lastNames)
	case "email":
		return fmt.Sprintf("%s%d@%s", strings.ToLower(randPick(firstNames)), randInt(999), randPick(domains))
	case "number", "age":
		return fmt.Sprintf("%d", 18+randInt(62))
	case "date":
		d := time.Now().AddDate(0, 0, -randInt(3650))
		return d.Format("2006-01-02")
	case "bool":
		if randInt(2) == 0 {
			return "true"
		}
		return "false"
	case "city":
		return randPick(cities)
	case "country":
		countries := []string{"US", "UK", "CA", "AU", "DE", "FR", "JP", "BR", "IN", "MX"}
		return randPick(countries)
	case "phone":
		return fmt.Sprintf("(%03d) %03d-%04d", 200+randInt(800), 200+randInt(800), randInt(10000))
	case "uuid":
		b := make([]byte, 16)
		rand.Read(b)
		b[6] = (b[6] & 0x0f) | 0x40
		b[8] = (b[8] & 0x3f) | 0x80
		return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	default:
		return fmt.Sprintf("val_%d", randInt(1000))
	}
}

// --- JSON subcommand ---
var jsonCount int

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Generate sample JSON objects",
	Long: `Generate random JSON objects with id, name, email, age, active, created_at.

Examples:
  openGyver testdata json --count 5`,
	RunE: func(c *cobra.Command, args []string) error {
		type Record struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Age       int    `json:"age"`
			Active    bool   `json:"active"`
			CreatedAt string `json:"created_at"`
		}
		var records []Record
		for i := 0; i < jsonCount; i++ {
			first := randPick(firstNames)
			last := randPick(lastNames)
			records = append(records, Record{
				ID:        i + 1,
				Name:      first + " " + last,
				Email:     fmt.Sprintf("%s.%s@%s", strings.ToLower(first), strings.ToLower(last), randPick(domains)),
				Age:       18 + randInt(62),
				Active:    randInt(2) == 1,
				CreatedAt: time.Now().AddDate(0, 0, -randInt(365)).Format(time.RFC3339),
			})
		}
		return cmd.PrintJSON(records)
	},
}

// --- Number subcommand ---
var (
	numMin, numMax, numCount int
	numFloat                 bool
)

var numberCmd = &cobra.Command{
	Use:   "number",
	Short: "Generate random numbers",
	Long: `Generate random numbers within a range.

Examples:
  openGyver testdata number --min 1 --max 100 --count 5
  openGyver testdata number --min 0 --max 1 --float --count 3`,
	RunE: func(c *cobra.Command, args []string) error {
		var values []interface{}
		for i := 0; i < numCount; i++ {
			if numFloat {
				v := float64(numMin) + float64(randInt(10000))/10000.0*float64(numMax-numMin)
				values = append(values, v)
				if !jsonOut {
					fmt.Printf("%.4f\n", v)
				}
			} else {
				v := numMin + randInt(numMax-numMin+1)
				values = append(values, v)
				if !jsonOut {
					fmt.Println(v)
				}
			}
		}
		if jsonOut {
			return cmd.PrintJSON(map[string]interface{}{"min": numMin, "max": numMax, "count": numCount, "values": values})
		}
		return nil
	},
}

func init() {
	testdataCmd.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "output as JSON")

	personCmd.Flags().IntVar(&personCount, "count", 1, "number of people to generate")
	csvCmd.Flags().IntVar(&csvRows, "rows", 10, "number of rows")
	csvCmd.Flags().StringVar(&csvColumns, "columns", "name,email,age,city", "comma-separated column types")
	jsonCmd.Flags().IntVar(&jsonCount, "count", 5, "number of records")
	numberCmd.Flags().IntVar(&numMin, "min", 0, "minimum value")
	numberCmd.Flags().IntVar(&numMax, "max", 100, "maximum value")
	numberCmd.Flags().IntVar(&numCount, "count", 1, "how many numbers")
	numberCmd.Flags().BoolVar(&numFloat, "float", false, "generate floating point numbers")

	testdataCmd.AddCommand(personCmd, csvCmd, jsonCmd, numberCmd)
	cmd.Register(testdataCmd)
}
