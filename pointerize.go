package stream

/* These helpers return pointers to provided base types;
 * these are needed because Go does not allow pointers to constant
 * expressions, or for the allocation of a pointer to base type to
 * in the same expression as its initialization.
 */

// IntPtr returns a pointer to an int.
func IntPtr(v int) *int { return &v }

// BoolPtr returns a pointer to a bool.
func BoolPtr(v bool) *bool { return &v }
