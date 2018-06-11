package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/djavorszky/depser"
)

var sources []string

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please specify one or more paths to check, or a filename with -f")
		os.Exit(1)
	}
	epoch := time.Now()

	fileName := flag.String("f", "none", "file to read for sources")

	flag.Parse()

	var err error
	if *fileName == "none" {
		sources = os.Args[1:]
	} else {
		sources, err = parseFile(*fileName)
		if err != nil {
			fmt.Printf("Failed parsing file: %v", err)
			os.Exit(1)
		}
	}

	start := time.Now()
	fmt.Println("Building dependencies")
	dep, err := depser.BuildDependencies(true, sources)
	if err != nil {
		fmt.Printf("failed building dependencies: %v\n", err)
	}
	fmt.Printf("Dependencies built in %s\n", time.Since(start))

	start = time.Now()
	fmt.Println("Checking cyclic dependencies")
	cyclics, ok := dep.CheckCyclicDependencies()
	if !ok {
		fmt.Printf("%d dependency cycle(s) detected:\n", len(cyclics))
		for _, cycle := range cyclics {
			fmt.Println(cycle)
		}
	}

	fmt.Printf("Cyclic dependency check done in %s\n", time.Since(start))
	fmt.Printf("Whole process took %s\n", time.Since(epoch))
}

func parseFile(fileName string) ([]string, error) {
	var sources []string

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Failed opening file: %v", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sources = append(sources, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading from io.Reader: %v", err)
	}

	return sources, nil
}
