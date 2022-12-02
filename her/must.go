package her

// Must panics if err is not nil and returns t otherwise.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// Must1 panics if err is not nil.
func Must1(err error) {
	if err != nil {
		panic(err)
	}
}
