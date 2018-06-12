package dependency

import (
	"fmt"
	"sort"
	"strings"
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

	knownCyclers sync.Map
	cyclicDeps   sync.Map
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

	if d.allowCycles {
		return nil
	}

	return d.checkCycles(depender)
}

// CheckCyclicDependencies checks to see if there are any cyclic dependencies.
// If there are, it will return them as a slice of strings.
//
// Boolean is set to true if no cyclic dependencies are found, and false
// if there are at least 1.
func (d *Dependency) CheckCyclicDependencies() ([]string, bool) {
	var sortableDeps []string
	for root := range d.deps {
		sortableDeps = append(sortableDeps, root)
	}

	sort.Strings(sortableDeps)

	var wg sync.WaitGroup

	wg.Add(len(sortableDeps))
	for _, root := range sortableDeps {
		d.checkCyclesAsync(root, &wg)
	}

	wg.Wait()

	var cycles []string

	d.cyclicDeps.Range(func(key, value interface{}) bool {
		cycles = append(cycles, key.(string))

		return true
	})

	return cycles, len(cycles) == 0
}

// checkCycles checks if there are any dependency cycles starting
// from depender
func (d *Dependency) checkCycles(depender string) error {
	seen := make(map[string]struct{})

	return d.check(depender, depender, seen)
}

// checkCyclesAsync checks if there are any dependency cycles starting
// from depender. It is meant to run in an asynchronous manner.
func (d *Dependency) checkCyclesAsync(depender string, wg *sync.WaitGroup) {
	seen := make(map[string]struct{})

	err := d.check(depender, depender, seen)
	if err != nil {
		d.cyclicDeps.Store(err.Error(), struct{}{})
	}

	wg.Done()
}

func (d *Dependency) check(route, depender string, seen map[string]struct{}) error {
	d.depRW.RLock()
	dependents := d.deps[depender]
	d.depRW.RUnlock()

	sort.Strings(dependents)
	for _, dep := range dependents {
		if _, ok := d.knownCyclers.Load(dep); ok {
			continue
		}

		depRoute := fmt.Sprintf("%s -> %s", route, dep)

		if _, ok := seen[dep]; ok {

			// double check
			if _, loaded := d.knownCyclers.LoadOrStore(dep, struct{}{}); !loaded {
				return fmt.Errorf(mustTrimToCycle(depRoute, dep))
			}

			continue
		}

		seen[dep] = struct{}{}

		err := d.check(depRoute, dep, seen)
		if err != nil {
			return err
		}
	}

	delete(seen, depender)

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

func trimToCycle(cycle, offender string) (string, error) {
	if cycle == "" || offender == "" {
		return "", fmt.Errorf("cycle or offender is empty")
	}

	ind := strings.Index(cycle, offender)
	if ind == -1 {
		return "", fmt.Errorf("not found in cycle: %s", offender)
	}

	return cycle[ind:], nil
}

func mustTrimToCycle(cycle, offender string) string {
	res, err := trimToCycle(cycle, offender)
	if err != nil {
		panic(err)
	}

	return res
}
