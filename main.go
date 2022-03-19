package main

import (
	"git-projects/cmd"
	"os"
)

func main() {

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}

}
