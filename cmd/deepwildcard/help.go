package main

import (
	"deepwildcard/internal/deepwildcard/version"
	"fmt"
)

func printUsage() {
	fmt.Println("Usage: deepwildcard [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --config file   path to config file.")
	fmt.Println("                    Default: ", defaultConfigFile)
	fmt.Println("  --log-time      Write time to output log")
	fmt.Println()
}

func printHeader() {
	fmt.Println()
	fmt.Printf("deepwildcard, a webhook microservice for step-ca\n")
	fmt.Printf("    version: %s\n", version.DW_VERSION)
	fmt.Printf("     github: %s\n", version.DW_URL)
	fmt.Println()
}
