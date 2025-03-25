package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/clickhouse-demo/internal"
)

func main() {
	// Parse command-line arguments
	example := flag.String("example", "sim", "Which example to run: 'sim' for simulation or 'real' for real implementation")
	help := flag.Bool("help", false, "Show help information")
	flag.Parse()

	if *help {
		fmt.Println("ClickHouse Exporter Demo for collex")
		fmt.Println("\nOptions:")
		fmt.Println("  -example=sim  Run the simulation example (default)")
		fmt.Println("  -example=real Run the real implementation example")
		fmt.Println("  -help         Show this help message")
		os.Exit(0)
	}

	// Run the selected example
	switch *example {
	case "sim":
		fmt.Println("Running simulated example...")
		internal.SimulatedClickHouseExample()
	case "real":
		fmt.Println("Running real implementation example...")
		internal.RealClickHouseExample()
	default:
		fmt.Printf("Unknown example: %s\n", *example)
		fmt.Println("Use -help for usage information")
		os.Exit(1)
	}
} 