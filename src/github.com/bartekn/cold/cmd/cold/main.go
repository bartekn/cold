package main

import (
	"fmt"
	"os"

	"github.com/bartekn/cold"
)

func main() {
	if err := cold.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
