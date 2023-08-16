package mysplit

import (
	"reflect"
	"testing"
)

func Test__shuffleRoundRobin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    []int
		n    int
		want [][]int
	}{
		{
			name: "1, 2, 3, 4, 5を2個ずつに分ける",
			a:    []int{1, 2, 3, 4, 5},
			n:    2,
			want: [][]int{
				{1, 3, 5},
				{2, 4},
			},
		},
		{
			name: "1, 2, 3, 4, 5, 6, 7を2個ずつに分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7},
			n:    2,
			want: [][]int{
				{1, 3, 5, 7},
				{2, 4, 6},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := roundRobin(test.a, test.n)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("got %v, want %v", test.a, test.want)
			}
		})
	}
}
