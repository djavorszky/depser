package depser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Hello is the entrypoint
func Hello() {
	srcs := []string{
		"/Users/javdaniel/git/liferay-support-ee/tools/environment-setup/LiferayUp/src/com",
		"/Users/javdaniel/git/liferay-support-ee/tools/environment-setup/LiferayUp/src/hu",
	}

	var wg sync.WaitGroup

	wg.Add(len(srcs))
	for _, src := range srcs {
		go walkPath(src, walker, &wg)
	}

	wg.Wait()
}

func walkPath(path string, walker filepath.WalkFunc, wg *sync.WaitGroup) {
	err := filepath.Walk(path, walker)
	if err != nil {
		log.Printf("error walking the path %q: %v", path, err)
	}

	wg.Done()
}

func walker(path string, info os.FileInfo, err error) error {
	const op = "walker"

	if err != nil {
		return fmt.Errorf("%v: can't visit %q: %v", op, path, err)
	}

	if info.IsDir() {
		log.Printf("visited folder: %q", path)
		return nil
	}

	if !strings.HasSuffix(info.Name(), ".java") {
		log.Printf("visited non-java file: %v", info.Name())
		return nil
	}

	imports, err := extractImports(path)
	if err != nil {
		return fmt.Errorf("%v: extract failed: %v", op, err)
	}

	for _, imp := range imports {
		log.Printf("Import: %q", imp)
	}

	return nil
}
