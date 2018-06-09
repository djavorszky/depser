package dependency

import (
	"testing"
)

func TestDependency_doAddDependency(t *testing.T) {
	type args struct {
		dependent string
		dependee  string
	}
	tests := []struct {
		name string
		args []args
	}{
		{"First", []args{args{"First", "dependee"}}},
		{"Two runs", []args{args{"Two runs", "dependee1"}, args{"Two runs", "dependee2"}}},
		{"Three runs", []args{args{"Three runs", "dependee1"}, args{"Three runs", "dependee2"}, args{"Three runs", "dependee3"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			for _, depTest := range tt.args {
				d.mustAddDependency(depTest.dependent, depTest.dependee)

				res, ok := d.deps[depTest.dependent]
				if !ok {
					t.Errorf("failed to add dependency: %s -> %s", depTest.dependent, depTest.dependee)
				}

				found := false
				for _, d := range res {
					if d == depTest.dependee {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("didn't find added dependency: %s -> %s", depTest.dependent, depTest.dependee)
				}
			}

			actualLen := len(d.deps[tt.name])
			if actualLen != len(tt.args) {
				t.Errorf("length of resulting dependency list is not as expected. Expected: %d, actual: %d", len(tt.args), actualLen)
			}
		})
	}
}

func TestDependency_Add(t *testing.T) {
	type args struct {
		dependent string
		dependee  string
	}
	tests := []struct {
		name    string
		args    []args
		wantErr bool
	}{
		{"First", []args{args{"First", "dependee"}}, false},
		{"Two runs", []args{args{"Two runs", "dependee1"}, args{"Two runs", "dependee2"}}, false},
		{"Three runs", []args{args{"Three runs", "dependee1"}, args{"Three runs", "dependee2"}, args{"Three runs", "dependee3"}}, false},

		{"Missing dependent", []args{args{"", "dependee"}}, true},
		{"Missing dependee", []args{args{"First", ""}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			for _, depTest := range tt.args {
				if err := d.Add(depTest.dependent, depTest.dependee); (err != nil) != tt.wantErr {
					t.Errorf("Dependency.Add() error = %v, wantErr %v", err, tt.wantErr)
				}

				if tt.wantErr {
					return
				}

				res, ok := d.deps[depTest.dependent]
				if !ok {
					t.Errorf("failed to add dependency: %s -> %s", depTest.dependent, depTest.dependee)
				}

				found := false
				for _, d := range res {
					if d == depTest.dependee {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("didn't find added dependency: %s -> %s", depTest.dependent, depTest.dependee)
				}
			}

			actualLen := len(d.deps[tt.name])
			if len(tt.args) != actualLen {
				t.Errorf("length of resulting dependency list is not as expected. Expected: %d, actual: %d", len(tt.args), actualLen)
			}
		})
	}
}
