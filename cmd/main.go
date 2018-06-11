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

	dep, err := depser.BuildDependencies(true, os.Args[1:])
	if err != nil {
		fmt.Printf("failed building dependencies: %v\n", err)
	}

	cyclics, ok := dep.CheckCyclicDependencies()
	if !ok {
		fmt.Printf("%d dependency cycle(s) detected:\n", len(cyclics))
		for _, cycle := range cyclics {
			fmt.Println(cycle)
		}
	}
}
