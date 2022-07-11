package self

import (
	"sync"
)

var (
	clearFuncs []func() error
	mut        sync.Mutex
)

func RegisterClearFuncs(f func() error) {
	mut.Lock()
	defer mut.Unlock()

	if f != nil {
		clearFuncs = append(clearFuncs, f)
	}
}

func RunClearFuncs() []error {
	mut.Lock()
	defer mut.Unlock()

	var e = make([]error, 0, len(clearFuncs))
	for i := len(clearFuncs) - 1; i >= 0; i-- {
		f := clearFuncs[i]
		if f != nil {
			if err := f(); err != nil {
				e = append(e, err)
			}
		}
	}
	return e
}
