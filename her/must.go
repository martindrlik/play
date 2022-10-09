package her

// Must is a helper that wraps a call to a function returning (T, error)
// and panics if the error is non-nil.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
