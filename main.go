package main

import (
	"os"

	"github.com/mikeh-elastic/nozzlebeat/cmd"
	_ "github.com/mikeh-elastic/nozzlebeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
