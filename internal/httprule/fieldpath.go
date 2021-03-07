package httprule

import "strings"

// FieldPath describes the path for a field from a message.
// Individual segments are in snake case (same as in protobuf file).
type FieldPath []string

func (f FieldPath) String() string {
	return strings.Join(f, ".")
}
