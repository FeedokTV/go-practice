package errorsgenerics

import (
	"cmp"
	"fmt"
)

// ----- Errors -----

type NotFoundError struct {
	ID string
}

func (e NotFoundError) Error() string { return fmt.Sprintf("id %s not found", e.ID) }

var _ error = NotFoundError{}

// ----- Generics -----

// MaxBy finds object with the biggest Orderable value (cmp.Ordered)
// It accepts slice of objects (T) and function that returns orderable
// field of the value. For the cases where two keys are equal - first seen
// will be chosen. For zero len slice second value returns false
func MaxBy[T any, K cmp.Ordered](ts []T, field func(t T) K) (T, bool) {
	var max T

	if len(ts) == 0 {
		return max, false
	}

	max = ts[0]

	for _, t := range ts[1:] {
		a := field(t)
		// For equaled keys - returns first max
		if a > field(max) {
			max = t
		}
	}

	return max, true
}
