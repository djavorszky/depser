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
	var (
		dep sync.RWMutex
		vis sync.RWMutex
	)

	dependency := Dependency{
		allowCycles:  allowCycles,
		deps:         make(map[string][]string),
		visibilities: make(map[string][]string),
		depRW:        &dep,
		visRW:        &vis,
	}

	return &dependency
}

// Add adds a new dependency to the depender, as well as set the
// corresponding visibility.
//
// If A depends on B, then A is the depender, and B is the dependent.
// In the above scenario, B needs to be visible to A.
func (d *Dependency) Add(depender, dependent string) error {
	if depender == "" || dependent == "" {
		return fmt.Errorf("empty dependant or dependee")
	}

	d.depRW.Lock()
	d.mustAddDependency(depender, dependent)
	d.depRW.Unlock()

	d.visRW.Lock()
	d.mustAddVisibility(dependent, depender)
	d.visRW.Unlock()

	return nil
}

// mustAddDependency is not concurrent-safe.
//
// If A depends on B, then A is the depender, and B is the dependent.
func (d *Dependency) mustAddDependency(depender, dependent string) {
	if depender == "" || dependent == "" {
		panic("empty dependent or dependee")
	}

	dependees := d.deps[depender]

	// Check if dependency already exists
	for _, dep := range dependees {
		if dependent == dep {
			return
		}
	}

	dependees = append(dependees, dependent)

	d.deps[depender] = dependees
}

// mustAddVisibility is not concurrent-safe.
//
// If A depends on B, then A is the depender, B is the dependent.
// In this scenario, B needs to be made visible to A, so we can say
// that B is being stalked by A, who is the stalker.
func (d *Dependency) mustAddVisibility(stalked, stalker string) {
	if stalked == "" || stalker == "" {
		panic("empty stalked or stalker")
	}

	stalkers := d.visibilities[stalked]

	// Check if stalking already exists
	for _, s := range stalkers {
		if stalker == s {
			return
		}
	}

	stalkers = append(stalkers, stalker)

	d.visibilities[stalked] = stalkers
}
