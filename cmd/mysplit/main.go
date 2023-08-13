package main

import (
	"flag"
	"log"
	"os"

	mysplit "github.com/matumoto1234/my-split"
	"github.com/matumoto1234/my-split/option"
)

func main() {
	opts, err := option.Parse()
	if err != nil {
		log.Fatal(err)
	}

	name := flag.Arg(0)
	prefix := flag.Arg(1)

	if prefix == "" {
		prefix = "x"
	}

	cli := &mysplit.CLI{
		Stdin: os.Stdin,
		Dir:   ".",
	}

	if err := cli.Run(name, prefix, opts); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
