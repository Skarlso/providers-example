package main

import (
	"fmt"
	"os"

	"github.com/Skarlso/providers-example/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("failed to run executor: ", err)
		os.Exit(1)
	}
}
