package depser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/djavorszky/depser/dependency"
)

var (
	dep  *dependency.Dependency
	errs []error
)

// BuildDependencies walks through all of the paths to build up a dependency tree
func BuildDependencies(allowCycles bool, roots []string) (*dependency.Dependency, error) {
	dep = dependency.NewWithCycles(allowCycles)

	var wg sync.WaitGroup

	wg.Add(len(roots))
	for _, root := range roots {
		go walkPath(root, walker, &wg)
	}

	wg.Wait()

	if len(errs) != 0 {
		var errMsg string
		for _, err := range errs {
			errMsg = fmt.Sprintf("%s && %v", errMsg, err)
		}

		errMsg = strings.TrimPrefix(errMsg, " && ")

		return nil, fmt.Errorf(errMsg)
	}

	return dep, nil
}

func walkPath(path string, walker filepath.WalkFunc, wg *sync.WaitGroup) {
	err := filepath.Walk(path, walker)
	if err != nil {
		errs = append(errs, fmt.Errorf("path %q: %v", path, err))
	}

	wg.Done()
}

func walker(path string, info os.FileInfo, err error) error {
	const op = "walker"

	if err != nil {
		return fmt.Errorf("%v: can't visit %q: %v", op, path, err)
	}

	if info.IsDir() {
		//log.Printf("visited folder: %q", path)
		return nil
	}

	if !strings.HasSuffix(info.Name(), ".java") {
		//log.Printf("visited non-java file: %v", info.Name())
		return nil
	}

	fqcn, err := extractFQCN(path)
	if err != nil {
		return fmt.Errorf("FQCN extract: %v", err)
	}

	imports, err := extractImports(path)
	if err != nil {
		return fmt.Errorf("%v: extract imports failed: %v", op, err)
	}

	for _, imp := range imports {
		err := dep.Add(fqcn, imp)
		if err != nil {
			return fmt.Errorf("failed adding dependency: %v", err)
		}
	}

	return nil
}
