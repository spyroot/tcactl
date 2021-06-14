package app

import (
	"strings"
)

const (
	StateInstantiate = "INSTANTIATE"
	StateCompleted   = "COMPLETED"
)

// IsInState - abstract state change,  late if I move code to
// state as different representation we swap code here.
func IsInState(currentState string, predicate string) bool {
	return strings.Contains(currentState, predicate)
}
