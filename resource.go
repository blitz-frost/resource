// Package resource provides manual reference tracking of arbitrary Go values.
//
// A notable use case for this package is working with variables that traverse the barrier between Go and another language.
package resource

import (
	"sync"
)

// A Pointer tracks a single arbitrary value. They should be treated as C pointers.
// In particular, any Pointer that is no longer needed should always be wiped.
type Pointer uint

var (
	live = make(map[Pointer]any) // currently reachable values
	next Pointer
	mux  sync.RWMutex
)

// Alloc returns a new Pointer that initially points at the given value.
func Alloc(value interface{}) Pointer {
	mux.Lock()
	defer mux.Unlock()

	// find free reference
	for {
		if _, ok := live[next]; !ok {
			break
		}
		next++
	}

	live[next] = value
	return next
}

// Get deferences the pointer.
// Returns nil if the pointer is invalid.
func (x Pointer) Get() interface{} {
	mux.RLock()
	v := live[x]
	mux.RUnlock()
	return v
}

// Set replaces the value currently referenced.
func (x Pointer) Set(value interface{}) {
	mux.Lock()
	live[x] = value
	mux.Unlock()
}

func (x Pointer) Wipe() {
	mux.Lock()
	delete(live, x)
	mux.Unlock()
}
