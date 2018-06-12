package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/djavorszky/depser"
)

var sources []string

func main() {
	if len(os.Args) == 1 {
		log.Println("Please specify one or more paths to check, or a filename with -f")
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
			log.Printf("Failed parsing file: %v", err)
			os.Exit(1)
		}
	}

	log.Printf("Got %d source files\n", len(sources))
	log.Println("Building dependencies")

	start := time.Now()

	dep, err := depser.BuildDependencies(true, sources)
	if err != nil {
		log.Printf("failed building dependencies: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Dependencies built in %s\n", time.Since(start))
	log.Println("Checking cyclic dependencies")

	start = time.Now()
	cyclics, ok := dep.CheckCyclicDependencies()
	if !ok {
		log.Printf("%d dependency cycle(s) detected:\n", len(cyclics))

		for _, cycle := range cyclics {
			log.Println(cycle)
		}

	}

	log.Printf("Cyclic dependency check done in %s\n", time.Since(start))
	log.Printf("Whole process took %s\n", time.Since(epoch))
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
