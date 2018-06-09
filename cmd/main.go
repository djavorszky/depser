package main

import (
	"fmt"
	"os"

	"github.com/djavorszky/depser"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please specify one or more paths to check")
		os.Exit(1)
	}

	err := depser.BuildDependencies(os.Args[1:])
	if err != nil {
		fmt.Printf("failed building dependencies: %v\n", err)
	}

}
