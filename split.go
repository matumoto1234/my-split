package mysplit

// N個ずつの要素に分ける
// e.g.
//
//	[1, 2, 3, 4, 5], n = 2 -> [[1, 2], [3, 4], [5]]
//	[1, 2, 3, 4, 5], n = 3 -> [[1, 2, 3], [4, 5]]
//	[[a], [bb], [ccc], [d], [ee]], n = 2 -> [[[a], [bbb]], [[ccc], [d]], [[ee]]]
func splitByN[T any](a []T, n int) [][]T {
	var b [][]T

	for i := 0; i < len(a); i += n {
		begin := i
		end := i + n

		if end > len(a) {
			end = len(a)
		}

		b = append(b, a[begin:end])
	}

	return b
}

// N個の要素に分ける
// e.g.
//
//	[1, 2, 3, 4, 5], n = 2 -> [[1, 2, 3], [4, 5]]
//	[1, 2, 3, 4, 5], n = 3 -> [[1, 2], [3, 4], [5]]
//	[1, 2, 3, 4, 5, 6, 7], n = 3 -> [[1, 2, 3], [4, 5], [6, 7]]
func splitN[T any](a []T, n int) [][]T {
	var b [][]T

	block := len(a) / n

	for i := 0; i < len(a); {
		if len(b) < len(a)%n {
			b = append(b, a[i:i+block+1])
			i += block + 1
		} else {
			b = append(b, a[i:i+block])
			i += block
		}
	}

	return b
}
