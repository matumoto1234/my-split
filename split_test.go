package mysplit

import (
	"reflect"
	"testing"
)

func Test__splitByN(t *testing.T) {
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
				{1, 2},
				{3, 4},
				{5},
			},
		},
		{
			name: "1, 2, 3, 4, 5, 6, 7を2個ずつに分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7},
			n:    2,
			want: [][]int{
				{1, 2},
				{3, 4},
				{5, 6},
				{7},
			},
		},
		{
			name: "1, 2, 3, 4, 5, 6, 7, 8を2個ずつに分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7, 8},
			n:    2,
			want: [][]int{
				{1, 2},
				{3, 4},
				{5, 6},
				{7, 8},
			},
		},
		{
			name: "1, 2, 3, 4, 5, 6, 7, 8, 9を4個ずつに分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			n:    4,
			want: [][]int{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := splitByN(test.a, test.n)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("want %v, but got %v", test.want, got)
			}
		})
	}
}

func Test__splitN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a    []int
		n    int
		want [][]int
	}{
		{
			name: "1, 2, 3, 4, 5, 6, 7を3個に分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7},
			n:    3,
			want: [][]int{
				{1, 2, 3},
				{4, 5},
				{6, 7},
			},
		},
		{
			name: "1, 2, 3, 4, 5, 6, 7, 8を3個に分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7, 8},
			n:    3,
			want: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8},
			},
		},
		{
			name: "1, 2, 3, 4, 5, 6, 7, 8, 9を3個に分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			n:    3,
			want: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
		},
		{
			name: "1, 2, 3, 4, 5, 6, 7, 8, 9, 10を2個に分ける",
			a:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			n:    2,
			want: [][]int{
				{1, 2, 3, 4, 5},
				{6, 7, 8, 9, 10},
			},
		},
		{
			name: "1, 2, 3, 4を3個に分ける",
			a:    []int{1, 2, 3, 4},
			n:    3,
			want: [][]int{
				{1, 2},
				{3},
				{4},
			},
		},
		{
			name: "1, 2, 3, 4, 5を3個に分ける",
			a:    []int{1, 2, 3, 4, 5},
			n:    3,
			want: [][]int{
				{1, 2},
				{3, 4},
				{5},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := splitN(test.a, test.n)

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("want %v, but got %v", test.want, got)
			}
		})
	}
}
