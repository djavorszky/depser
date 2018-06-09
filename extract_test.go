package depser

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

const (
	validImports           = "import com.example.Test;\nimport com.example.util.Something;\nimport hu.coolio.Reader;"
	validImportsWithPrefix = "package Something;\n" + validImports

	invalidImports = "impor hello.something"

	validPackage           = "package Something;\n"
	validPackageWithSuffix = validPackage + validImports

	invalidPackage = "packge hello.something"
)

func Test_extractImportFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"valid", args{strings.NewReader(validImports)},
			[]string{"com.example.Test", "com.example.util.Something", "hu.coolio.Reader"}, false},
		{"valid with prefix", args{strings.NewReader(validImportsWithPrefix)},
			[]string{"com.example.Test", "com.example.util.Something", "hu.coolio.Reader"}, false},

		{"invalid", args{strings.NewReader(invalidImports)}, []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractImportFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractImportFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractImportFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractPackageFrom(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"valid", args{strings.NewReader(validPackage)}, "Something", false},
		{"valid with suffix", args{strings.NewReader(validPackageWithSuffix)}, "Something", false},

		{"invalid", args{strings.NewReader(invalidPackage)}, "", false}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractPackageFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractPackageFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractPackageFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
