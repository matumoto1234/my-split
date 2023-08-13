package option

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
)

type SplitWay string

const (
	ByLine SplitWay = "l"
	ByByte SplitWay = "b"
)

type Option struct {
	Way SplitWay
	L   int
	B   int64
}

func Parse() (*Option, error) {
	var (
		l = flag.Int("l", 1000, "put NUMBER lines/records per output file")
		b = flag.String("b", "3MB", "put NUMBER bytes per output file")
	)

	flag.Parse()

	o := &Option{}
	var err error

	o.Way, err = parseWay()
	if err != nil {
		return nil, err
	}

	o.L, err = parseL(l)
	if err != nil {
		return nil, err
	}

	o.B, err = parseB(b)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func specified(name string) bool {
	is := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			is = true
		}
	})

	return is
}

func parseWay() (SplitWay, error) {
	if specified("l") && specified("b") {
		return ByLine, fmt.Errorf("you may only specify one of the -l and -b options")
	}

	if specified("l") {
		return ByLine, nil
	}

	if specified("b") {
		return ByByte, nil
	}

	return ByLine, nil
}

func parseL(l *int) (int, error) {
	if !specified("l") && specified("b") {
		// bのみ指定された場合のデフォルト値
		// 1 << 25 は適当に決定したそこそこ大きな値
		return 1 << 25, nil
	}

	if *l < 0 {
		return 0, fmt.Errorf("invalid number of lines: '%d'", *l)
	}

	return *l, nil
}

func parseB(b *string) (int64, error) {
	if specified("l") && !specified("b") {
		// lのみ指定された場合のデフォルト値
		// 1 << 25 は適当に決定したそこそこ大きな値
		return 1 << 25, nil
	}

	// match example : 3, 1gb, 1b, 30GB, 27g, 1000kb
	reg := regexp.MustCompile(`^(\d+)(.+)?$`)

	matched := reg.FindSubmatch([]byte(*b))
	if matched == nil {
		return 0, fmt.Errorf("invalid number of bytes: '%s'", *b)
	}

	size, err := strconv.ParseInt(string(matched[1]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number of bytes: '%s'", *b)
	}

	if size < 0 {
		return 0, fmt.Errorf("invalid number of bytes: '%s'", *b)
	}

	unit := string(matched[2])

	switch unit {
	case "":
		// pass
	case "b":
		size *= 512
	case "k", "K", "KiB":
		size *= 1024
	case "m", "M", "MiB":
		size *= 1024 * 1024
	case "g", "G", "GiB":
		size *= 1024 * 1024 * 1024
	case "kb", "KB":
		size *= 1000
	case "mb", "MB":
		size *= 1000 * 1000
	case "gb", "GB":
		size *= 1000 * 1000 * 1000
	default:
		return 0, fmt.Errorf("invalid number of bytes: '%s'", *b)
	}

	return size, nil
}
