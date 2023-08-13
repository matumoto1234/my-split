package mysplit

import "regexp"

// N個ずつの行に分割する
func splitByNLines(bytes []byte, n int) [][]string {
	lines := regexp.MustCompile("\r\n|\n").Split(string(bytes), -1)
	nLines := splitByN(lines, n)
	return nLines
}

// e.g. [a, b, c, d, e], n = 2 -> [[a, b], [c, d], [e]]
func splitByN(strs []string, n int) [][]string {
	var a [][]string

	for i := 0; i < len(strs); i += n {
		begin := i
		end := i + n

		if end > len(strs) {
			end = len(strs)
		}

		a = append(a, strs[begin:end])
	}

	return a
}
