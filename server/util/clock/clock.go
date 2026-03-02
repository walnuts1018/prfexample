package clock

import (
	"time"

	"github.com/Code-Hex/synchro"
)

type Clock[T synchro.TimeZone] struct {
	Now func() synchro.Time[T]
}

func Default[T synchro.TimeZone]() Clock[T] {
	return Clock[T]{
		Now: func() synchro.Time[T] {
			return synchro.Now[T]()
		},
	}
}

func Dummy[T synchro.TimeZone](override ...synchro.Time[T]) Clock[T] {
	var now func() synchro.Time[T]
	if len(override) > 0 {
		now = func() synchro.Time[T] {
			return override[0]
		}
	} else {
		now = func() synchro.Time[T] {
			return synchro.New[T](2006, time.January, 2, 15, 4, 5, 0)
		}
	}

	return Clock[T]{
		Now: now,
	}
}
