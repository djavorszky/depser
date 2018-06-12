package depser

import (
	"fmt"
	"strings"
)

func parseImport(line string) (string, error) {
	const op = "parseImport"

	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "import") {
		return "", fmt.Errorf("%v: import statement not found: %v", op, line)
	}

	// Line has to be at least 9 characters for it to be a valid import:
	// import a;
	if len(line) < 9 {
		return "", fmt.Errorf("%v: line too short to be valid: %s", op, line)
	}

	// Remove the "import a;" from the beginning
	// and the semicolon at the end.
	return line[7 : len(line)-1], nil
}

func mustParseImport(line string) string {
	imp, err := parseImport(line)
	if err != nil {
		panic("parse failed: " + err.Error())
	}

	return imp
}

func parsePackage(line string) (string, error) {
	const op = "parsePackage"

	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "package") {
		return "", fmt.Errorf("%v: package statement not found: %v", op, line)
	}

	// Line has to be at least 10 characters for it to be a valid package declaration:
	// package a;
	if len(line) < 10 {
		return "", fmt.Errorf("%v: line too short to be valid: %s", op, line)
	}

	// Remove the "package " from the beginning and the semicolon at the end.
	return line[8 : len(line)-1], nil
}

func mustParsePackage(line string) string {
	pkg, err := parsePackage(line)
	if err != nil {
		panic("parse failed: " + err.Error())
	}

	return pkg
}
