package convertfile

import "fmt"

// Apple Numbers files use the IWA (iWork Archive) format — a proprietary
// protobuf-based container inside a zip file. There is no mature Go library
// for reading or writing .numbers files natively.
//
// For now these conversions return a clear error directing users to use
// CSV or XLSX as an intermediary. When a suitable library becomes available,
// these stubs can be replaced with real implementations.

func csvToNumbers(opts ConvertOpts) error {
	return fmt.Errorf("CSV → Numbers conversion is not yet implemented.\n" +
		"Workaround: convert CSV → XLSX first, then open in Numbers and save as .numbers.\n" +
		"  openGyver convertFile data.csv -o data.xlsx")
}

func xlsxToNumbers(opts ConvertOpts) error {
	return fmt.Errorf("XLSX → Numbers conversion is not yet implemented.\n" +
		"Workaround: open the XLSX file in Apple Numbers and save as .numbers.")
}

func numbersToCSV(opts ConvertOpts) error {
	return fmt.Errorf("Numbers → CSV conversion is not yet implemented.\n" +
		"Workaround: open the .numbers file in Apple Numbers, export as CSV, then use that file.\n" +
		"  File → Export To → CSV")
}

func numbersToXLSX(opts ConvertOpts) error {
	return fmt.Errorf("Numbers → XLSX conversion is not yet implemented.\n" +
		"Workaround: open the .numbers file in Apple Numbers, export as XLSX.\n" +
		"  File → Export To → Excel")
}
