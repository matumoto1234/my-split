package mysplit

// See : https://www.gnu.org/software/coreutils/manual/html_node/split-invocation.html#split-invocation
func nextSuffix(s string) string {
	if s == "" {
		return "aa"
	}

	leadingZ := 0
	for i := 0; i < len(s); i++ {
		if s[i] == 'z' {
			leadingZ++
		} else {
			break
		}
	}

	s2 := []rune(s[leadingZ:])

	for i := len(s2) - 1; i >= 0; i-- {
		if s2[i] != 'z' {
			s2[i]++
			break
		} else {
			s2[i] = 'a'
		}
	}

	if s2[0] == 'z' {
		return s[:leadingZ] + string(s2) + "aa"
	} else {
		return s[:leadingZ] + string(s2)
	}
}
