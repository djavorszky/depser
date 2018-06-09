package dependency

import (
	"fmt"
	"sync"
)

// Dependency holds the whole dependency list. Which class
// depends on which classes (that needs to be imported), as
// well as which class are visible to which packages
type Dependency struct {
	depRW        *sync.RWMutex
	deps         map[string][]string
	visRW        *sync.RWMutex
	visibilities map[string][]string
	allowCycles  bool
}

// New returns a ready-to-use Dependency struct. By default it is
// set not to allow dependency cycles. See NewWithCycles if you
// want to enable it
func New() *Dependency {
	return NewWithCycles(false)
}

// NewWithCycles returns a Dependency struct in which you can
// decide whether to allow dependency cycles or not.
func NewWithCycles(allowCycles bool) *Dependency {
	dependency := &Dependency{
		allowCycles:  allowCycles,
		deps:         make(map[string][]string),
		visibilities: make(map[string][]string),
	}

	return dependency
}

// Add adds a new dependency to the dependent, as well
// as sets the corresponding visibility if needed.
func (d *Dependency) Add(dependent, dependee string) error {
	if dependent == "" || dependee == "" {
		return fmt.Errorf("empty dependant or dependee")
	}

	d.mustAddDependency(dependent, dependee)

	return nil
}

// mustAddDependency is not concurrent-safe.
func (d *Dependency) mustAddDependency(dependent, dependee string) {
	if dependent == "" || dependee == "" {
		panic("empty dependent or dependee")
	}

	dependees := d.deps[dependent]

	// Check if dependency already exists
	for _, dep := range dependees {
		if dependee == dep {
			return
		}
	}

	dependees = append(dependees, dependee)

	d.deps[dependent] = dependees
}
