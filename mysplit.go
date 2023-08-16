package mysplit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/matumoto1234/my-split/option"
	"github.com/pkg/errors"
)

type CLI struct {
	Stdin  io.Reader
	Stdout io.Writer
	Dir    string
}

func (cli *CLI) Run(name, prefix string, opts *option.Option) error {
	r, err := cli.open(name)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := cli.run(r, prefix, opts); err != nil {
		return err
	}

	return nil
}

func (cli *CLI) open(name string) (io.Reader, error) {
	if name == "" || name == "-" {
		return cli.Stdin, nil
	}

	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("cannot open '%s' for reading: %s", name, err)
	}

	buf := bufio.NewReader(file)

	return buf, nil
}

func (cli *CLI) run(r io.Reader, prefix string, opt *option.Option) error {
	var suffix string
	var chunkIndex int

	for {
		input := make([]byte, opt.Byte)
		n, err := r.Read(input)
		input = input[:n]

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch opt.SplitWay {
		case option.ByLine:
			onlyLF := regexp.MustCompile(`\r\n|\r|\n`).ReplaceAll(input, []byte("\n"))
			lines := bytes.Split(onlyLF, []byte("\n"))
			nLines := splitByN(lines, opt.Line)

			for _, nLine := range nLines {
				data := bytes.Join(nLine, []byte("\n"))

				suffix = nextSuffix(suffix)

				if err := cli.write(prefix, suffix, data); err != nil {
					return err
				}
			}
		case option.ByByte:
			suffix = nextSuffix(suffix)

			if err := cli.write(prefix, suffix, input); err != nil {
				return err
			}
		case option.ByChunk:
			// ここでクロージャーを使っているのは、suffixとchunkIndexを更新するため
			// TODO: もしうまい具合に関数に切り出せるのであれば切り出す
			writeContents := func(contents [][]byte) error {
				for _, content := range contents {
					suffix = nextSuffix(suffix)
					chunkIndex++

					if opt.Chunk.K == nil {
						if err := cli.write(prefix, suffix, content); err != nil {
							return err
						}
					} else {
						if chunkIndex == *opt.Chunk.K {
							fmt.Fprintln(cli.Stdout, string(content))
							return nil
						}
					}
				}

				return nil
			}

			switch opt.Chunk.Type {
			case option.LineChunk:
				onlyLF := regexp.MustCompile(`\r\n|\r|\n`).ReplaceAll(input, []byte("\n"))
				lines := bytes.Split(onlyLF, []byte("\n"))
				linesList := splitN(lines, opt.Chunk.N)

				contents := make([][]byte, len(linesList))

				for i, content := range linesList {
					joined := bytes.Join(content, []byte("\n"))
					contents[i] = joined
				}

				if err := writeContents(contents); err != nil {
					return err
				}

			case option.ByteChunk:
				contents := splitN(input, opt.Chunk.N)

				if err := writeContents(contents); err != nil {
					return err
				}
			case option.RoundRobinChunk:
				onlyLF := regexp.MustCompile(`\r\n|\r|\n`).ReplaceAll(input, []byte("\n"))
				lines := bytes.Split(onlyLF, []byte("\n"))
				linesList := splitN(lines, opt.Chunk.N)
				shuffleRoundRobin(linesList, opt.Chunk.N)

				contents := make([][]byte, len(linesList))

				for i, content := range linesList {
					joined := bytes.Join(content, []byte("\n"))
					contents[i] = joined
				}

				if err := writeContents(contents); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (cli *CLI) write(prefix, suffix string, data []byte) error {
	err := os.WriteFile(filepath.Join(cli.Dir, prefix+suffix), []byte(data), 0644)

	if err != nil {
		return errors.WithMessagef(err, "cannot write to file: %s", filepath.Join(cli.Dir, prefix+suffix))
	}

	return nil
}

func shuffleRoundRobin[T any](a [][]T, n int) {
	var idx int
	b := make([][]T, n)

	for i := range a {
		b[idx] = append(b[idx], a[i]...)

		idx++
		if idx >= n {
			idx = 0
		}
	}
}
