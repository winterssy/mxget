package main

import (
	"os"

	"github.com/winterssy/mxget/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
