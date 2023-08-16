package mysplit

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matumoto1234/my-split/option"
)

func Test_CLI_run(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		opts       *option.Option
		in         string
		fileToWant map[string]string
		wantErr    bool
	}{
		{
			name:   "4行のテキストを2行ごとに分割した場合、xaaとxabに分割される",
			prefix: "x",
			opts: &option.Option{
				SplitWay: option.ByLine,
				Line:     2,
				Byte:     1000,
			},
			in: "1\n2\n3\n4",
			fileToWant: map[string]string{
				"xaa": "1\n2",
				"xab": "3\n4",
			},
			wantErr: false,
		},
		{
			name:   "5行のテキストを2行ごとに分割した場合、xaaとxabとxacに分割される",
			prefix: "x",
			opts: &option.Option{
				SplitWay: option.ByLine,
				Line:     2,
				Byte:     1000,
			},
			in: "1\n2\n3\n4\n5",
			fileToWant: map[string]string{
				"xaa": "1\n2",
				"xab": "3\n4",
				"xac": "5",
			},
			wantErr: false,
		},
		{
			name:   "4行のテキストを2バイトごとに分割した場合、xaaとxabとxacとxadに分割される",
			prefix: "x",
			opts: &option.Option{
				SplitWay: option.ByByte,
				Line:     1000,
				Byte:     2,
			},
			in: "1\n2\n3\n4",
			fileToWant: map[string]string{
				"xaa": "1\n",
				"xab": "2\n",
				"xac": "3\n",
				"xad": "4",
			},
			wantErr: false,
		},
		{
			name:   "4行のテキストを2チャンクで分割した場合、xaaとxabに分割される",
			prefix: "x",
			opts: &option.Option{
				SplitWay: option.ByChunk,
				Line:     1000,
				Byte:     1000,
				Chunk: &option.Chunk{
					Type: option.ByteChunk,
					K:    nil,
					N:    2,
				},
			},
			in: "1\n2\n3\n4",
			fileToWant: map[string]string{
				"xaa": "1\n2\n",
				"xab": "3\n4",
			},
			wantErr: false,
		},
		{
			name:   "4行のテキストをl/3チャンクで分割した場合、xaaとxabとxacに分割される",
			prefix: "x",
			opts: &option.Option{
				SplitWay: option.ByChunk,
				Line:     1000,
				Byte:     1000,
				Chunk: &option.Chunk{
					Type: option.LineChunk,
					K:    nil,
					N:    3,
				},
			},
			in: "1\n2\n3\n4",
			fileToWant: map[string]string{
				"xaa": "1\n2\n",
				"xab": "3\n",
				"xac": "4",
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		test := test

		// t.TempDirに生成されたファイルの中身を確認してテストを行う
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			d := t.TempDir()

			cli := &CLI{
				Stdin: strings.NewReader(test.in),
				Dir:   d,
			}

			r := strings.NewReader(test.in)

			err := cli.run(r, test.prefix, test.opts)

			switch {
			case test.wantErr && err == nil:
				t.Fatal("expected error did not occur")
			case !test.wantErr && err != nil:
				t.Fatalf("unexpected error: %+v", err)
			}

			files, err := os.ReadDir(d)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			for _, f := range files {
				want, ok := test.fileToWant[f.Name()]
				if !ok {
					t.Fatalf("generated unexpected file: %s", f.Name())
				}

				got, err := os.ReadFile(filepath.Join(d, f.Name()))
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}

				if !bytes.Equal(got, []byte(want)) {
					t.Errorf("file name: %s, want: `%s`, got: `%s`", f.Name(), want, got)
				}
			}
		})
	}
}
