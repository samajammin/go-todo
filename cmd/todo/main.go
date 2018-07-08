package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/sbrichards/go-todo/todo"
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
	case "done":
		err = done(strings.Join(flag.Args()[1:], " "))
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

type length int64

const (
	sizeOfLength = 8
	dbPath       = "tododb.pb"
)

var endianness = binary.LittleEndian

func add(text string) error {
	if len(text) < 1 {
		return fmt.Errorf("cannot add empty task")
	}

	task := &todo.Task{
		Title: text,
		Done:  false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("could not encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return fmt.Errorf("could not encode length of message: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("could not close file %s: %v", dbPath, err)
	}

	fmt.Println("\nAdded task,", proto.MarshalTextString(task))
	return nil
}

func list() error {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	fmt.Println("\nHere are your todos:\n")
	for {
		if len(b) == 0 {
			return nil
		} else if len(b) < sizeOfLength {
			return fmt.Errorf("remaining odd %d bytes", len(b))
		}

		var l length
		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[sizeOfLength:]

		var task todo.Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]

		if task.Done {
			fmt.Printf("[X]: ")
		} else {
			fmt.Printf("[ ]: ")
		}
		fmt.Println(task.Title)
	}
}

func done(text string) error {
	fmt.Println("NEED TO IMPLEMENT DONE!")
	return nil
}
