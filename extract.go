package depser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func extractImports(path string) ([]string, error) {
	const op = "extractImports"

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%v: failed to open %q: %v", op, path, err)
	}
	defer file.Close()

	return extractImportFrom(file)
}

func extractImportFrom(r io.Reader) ([]string, error) {
	const op = "exportImportFrom(io.Reader)"

	imports := make([]string, 0)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "import ") {
			imports = append(imports, mustParseImport(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%v: reading from io.Reader: %v", op, err)
	}

	return imports, nil
}

func extractPackage(path string) (string, error) {
	const op = "extractPackage"

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("%v: failed to open %q: %v", op, path, err)
	}
	defer file.Close()

	return extractPackageFrom(file)
}

func extractPackageFrom(r io.Reader) (string, error) {
	const op = "exportPackageFrom(io.Reader)"

	var pkg string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "package") {
			pkg = mustParsePackage(line)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("%v: reading from io.Reader: %v", op, err)
	}

	return pkg, nil
}
