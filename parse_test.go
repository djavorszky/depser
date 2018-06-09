package depser

import (
	"testing"
)

func Test_parseImport(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"working", args{"import com.liferay.test;"}, "com.liferay.test", false},
		{"working one", args{"import com;"}, "com", false},
		{"working shortest", args{"import a;"}, "a", false},

		{"corner case", args{"import;"}, "", true},
		{"corner case2", args{"import ;"}, "", true},

		{"empty", args{""}, "", true},
		{"no import", args{"public class TestClass {"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseImport(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseImport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mustParseImport(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"working", args{"import com.liferay.test;"}, "com.liferay.test"},
		{"working one", args{"import com;"}, "com"},
		{"working shortest", args{"import a;"}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustParseImport(tt.args.line); got != tt.want {
				t.Errorf("mustParseImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePackage(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"working", args{"package com.liferay.test;"}, "com.liferay.test", false},
		{"working shorter", args{"package com.liferay;"}, "com.liferay", false},
		{"working one", args{"package com;"}, "com", false},
		{"working blob", args{"package java.io.*;"}, "java.io.*", false},
		{"working shortest", args{"package a;"}, "a", false},

		{"corner case", args{"package;"}, "", true},
		{"corner case2", args{"package;"}, "", true},
		{"empty", args{""}, "", true},
		{"no import", args{"public class TestClass {"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePackage(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mustParsePackage(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"working", args{"package com.liferay.test;"}, "com.liferay.test"},
		{"working shorter", args{"package com.liferay;"}, "com.liferay"},
		{"working one", args{"package com;"}, "com"},
		{"working blob", args{"package java.io.*;"}, "java.io.*"},
		{"working shortest", args{"package a;"}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustParsePackage(tt.args.line); got != tt.want {
				t.Errorf("mustParsePackage() = %v, want %v", got, tt.want)
			}
		})
	}
}
