package main

import (
	"flag"
	"log"
	"os"

	mysplit "github.com/matumoto1234/my-split"
	"github.com/matumoto1234/my-split/option"
)

func main() {
	mylog := log.New(os.Stderr, "mysplit: ", 0)

	opt, err := option.Parse()
	if err != nil {
		mylog.Fatal(err)
	}

	name := flag.Arg(0)
	prefix := flag.Arg(1)

	if prefix == "" {
		prefix = "x"
	}

	cli := &mysplit.CLI{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Dir:    ".",
	}

	if err := cli.Run(name, prefix, opt); err != nil {
		mylog.Fatal(err)
	}

	os.Exit(0)
}
