package uuid

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mj/opengyver/cmd"
	"github.com/spf13/cobra"
)

var (
	version   int
	count     int
	uppercase bool
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate UUID identifiers",
	Long: `Generate universally unique identifiers (UUIDs).

VERSIONS:

  v4   Random UUID (default). 122 bits of randomness. Best for most use cases.
  v6   Reordered time-based UUID. Sortable by creation time, includes a
       timestamp and random node. Good for database primary keys.

Use --count to generate multiple UUIDs at once.

Examples:
  openGyver uuid
  openGyver uuid --version 4
  openGyver uuid --version 6
  openGyver uuid --count 5
  openGyver uuid --version 6 --count 10
  openGyver uuid --uppercase`,
	Args: cobra.NoArgs,
	RunE: runUUID,
}

func runUUID(c *cobra.Command, args []string) error {
	for i := 0; i < count; i++ {
		id, err := generate(version)
		if err != nil {
			return err
		}
		s := id.String()
		if uppercase {
			s = strings.ToUpper(s)
		}
		fmt.Println(s)
	}
	return nil
}

func generate(ver int) (uuid.UUID, error) {
	switch ver {
	case 4:
		return uuid.NewRandom()
	case 6:
		return uuid.NewV6()
	default:
		return uuid.UUID{}, fmt.Errorf("unsupported UUID version: %d (supported: 4, 6)", ver)
	}
}

func init() {
	uuidCmd.Flags().IntVar(&version, "version", 4, "UUID version: 4 (random) or 6 (time-sorted)")
	uuidCmd.Flags().IntVar(&count, "count", 1, "number of UUIDs to generate")
	uuidCmd.Flags().BoolVar(&uppercase, "uppercase", false, "output in uppercase")
	cmd.Register(uuidCmd)
}
