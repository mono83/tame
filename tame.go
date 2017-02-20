package main

import (
	"fmt"
	"github.com/mono83/tame/cmd"
	"os"
)

func main() {
	if err := cmd.TameCmd.Execute(); err != nil {
		fmt.Println("Execution error occurred:", err)
		os.Exit(1)
	}
}
