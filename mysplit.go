package mysplit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/matumoto1234/my-split/option"
	"github.com/pkg/errors"
)

type CLI struct {
	Stdin  io.Reader
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

func (cli *CLI) run(r io.Reader, prefix string, opts *option.Option) error {
	var suf string

	for {
		body := make([]byte, opts.B)
		n, err := r.Read(body)
		body = body[:n]

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch option.SplitWay(opts.Way) {
		case option.ByLine:
			nLines := splitByNLines(body, opts.L)

			for _, nLine := range nLines {
				data := strings.Join(nLine, "\n")

				suf = nextSuffix(suf)

				if err := cli.write(prefix, suf, []byte(data)); err != nil {
					return err
				}
			}
		case option.ByByte:
			suf = nextSuffix(suf)

			if err := cli.write(prefix, suf, body); err != nil {
				return err
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
