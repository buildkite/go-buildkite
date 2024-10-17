package buildkite

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}

	return val
}
