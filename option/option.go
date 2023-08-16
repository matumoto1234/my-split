package option

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
)

type chunkType int

const (
	ByteChunk chunkType = iota
	LineChunk
	RoundRobinChunk
)

type Chunk struct {
	Type chunkType
	K    *int
	N    int
}

type splitWay int

const (
	ByLine splitWay = iota
	ByByte
	ByChunk
)

type Option struct {
	SplitWay splitWay
	Line     int
	Byte     int64
	Chunk    *Chunk
}

func Parse() (*Option, error) {
	var (
		l = flag.Int("l", 1000, "put NUMBER lines/records per output file")
		b = flag.String("b", "3MB", "put NUMBER bytes per output file")
		n = flag.String("n", "", "put NUMBER records per output file")
	)

	flag.Parse()

	o := &Option{}
	var err error

	o.SplitWay, err = parseSplitWay()
	if err != nil {
		return nil, err
	}

	o.Line, err = parseL(l)
	if err != nil {
		return nil, err
	}

	o.Byte, err = parseB(b)
	if err != nil {
		return nil, err
	}

	o.Chunk, err = parseN(n)
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

func parseSplitWay() (splitWay, error) {
	if specified("l") && specified("b") {
		return ByLine, fmt.Errorf("you may only specify one of the -l and -b options")
	}

	if specified("l") {
		return ByLine, nil
	}

	if specified("b") {
		return ByByte, nil
	}

	if specified("n") {
		return ByChunk, nil
	}

	return ByLine, nil
}

// parseL parses the -l option.
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

// parseB parses the -b option.
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

// parseN parses the -n option as below.
//   - n      generate n files based on current size of input
//   - k/n    output only kth of n to standard output
//   - l/n    generate n files without splitting lines or records
//   - l/k/n  likewise but output only kth of n to stdout
//   - r/n    like ‘l’ but use round robin distribution
//   - r/k/n  likewise but output only kth of n to stdout
//
// See also: https://www.gnu.org/software/coreutils/manual/html_node/split-invocation.html#split-invocation
func parseN(n *string) (*Chunk, error) {
	if !specified("n") {
		return nil, nil
	}

	// n
	onlyNum := regexp.MustCompile(`^(\d+)$`)
	if onlyNum.MatchString(*n) {
		num, err := strconv.Atoi(*n)
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		if num < 0 {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		return &Chunk{
			Type: ByteChunk,
			N:    num,
		}, nil
	}

	// k/n
	kNum := regexp.MustCompile(`^(\d+)/(\d+)$`)
	if kNum.MatchString(*n) {
		matched := kNum.FindStringSubmatch(*n)

		k, err := strconv.Atoi(matched[1])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		num, err := strconv.Atoi(matched[2])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		if k <= 0 || k > num {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		return &Chunk{
			Type: ByteChunk,
			N:    num,
			K:    &k,
		}, nil
	}

	// l/n
	lNum := regexp.MustCompile(`^l/(\d+)$`)
	if lNum.MatchString(*n) {
		matched := lNum.FindStringSubmatch(*n)

		num, err := strconv.Atoi(matched[1])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		return &Chunk{
			Type: LineChunk,
			N:    num,
		}, nil
	}

	// l/k/n
	lKNum := regexp.MustCompile(`^l/(\d+)/(\d+)$`)
	if lKNum.MatchString(*n) {
		matched := lKNum.FindStringSubmatch(*n)

		k, err := strconv.Atoi(matched[1])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		num, err := strconv.Atoi(matched[2])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		if k <= 0 || k > num {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		return &Chunk{
			Type: LineChunk,
			N:    num,
			K:    &k,
		}, nil
	}

	// r/n
	rNum := regexp.MustCompile(`^r/(\d+)$`)
	if rNum.MatchString(*n) {
		matched := rNum.FindStringSubmatch(*n)

		num, err := strconv.Atoi(matched[1])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		return &Chunk{
			Type: RoundRobinChunk,
			N:    num,
		}, nil
	}

	// r/k/n
	rKNum := regexp.MustCompile(`^r/(\d+)/(\d+)$`)
	if rKNum.MatchString(*n) {
		matched := rKNum.FindStringSubmatch(*n)

		k, err := strconv.Atoi(matched[1])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		if k <= 0 {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		num, err := strconv.Atoi(matched[2])
		if err != nil {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		if k <= 0 || k > num {
			return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
		}

		return &Chunk{
			Type: RoundRobinChunk,
			N:    num,
			K:    &k,
		}, nil
	}

	return nil, fmt.Errorf("invalid chunk number: '%s'", *n)
}
