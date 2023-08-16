package mysplit

func roundRobin[T any](a []T, n int) [][]T {
	b := make([][]T, n)

	for i := 0; i < len(a); i++ {
		b[i%n] = append(b[i%n], a[i])
	}

	return b
}
