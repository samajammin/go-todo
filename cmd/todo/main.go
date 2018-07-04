package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: add or list")
		os.Exit(1)
	}

	var err error
	switch cmd := flag.Arg(0); cmd {
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	case "list":
		err = list()
	default:
		err = fmt.Errorf("Unknown subcommand: %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func add(text string) error {
	fmt.Println("NEED TO IMPLEMENT ADD")
	return nil
}

func list() error {
	fmt.Println("NEED TO IMPLEMENT LIST")
	return nil
}
