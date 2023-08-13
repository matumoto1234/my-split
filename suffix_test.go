package mysplit

import "testing"

func Test__nextSuffix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: "aa"},
		{s: "aa", want: "ab"},
		{s: "ab", want: "ac"},
		{s: "yy", want: "yz"},
		{s: "yz", want: "zaaa"},
		{s: "zaaa", want: "zaab"},
		{s: "zyzz", want: "zzaaaa"},
		{s: "zzaaaa", want: "zzaaab"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			got := nextSuffix(tt.s)
			if got != tt.want {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}
