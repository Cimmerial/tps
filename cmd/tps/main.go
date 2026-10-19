package main

import (
	"os"

	"github.com/Cimmerial/tps/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}