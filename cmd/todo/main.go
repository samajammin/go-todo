package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sbrichards/go-todo/todo"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand: add, list, or done")
		os.Exit(1)
	}

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	client := todo.NewTasksClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "add":
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	case "done":
		err = done(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	case "list":
		err = list(context.Background(), client)
	default:
		err = fmt.Errorf("Unknown subcommand: %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func add(ctx context.Context, client todo.TasksClient, text string) error {
	t, err := client.Add(ctx, &todo.TaskTitle{Title: text})
	if err != nil {
		return fmt.Errorf("could not create task: %v", err)
	}
	fmt.Println("added task:", t.Title)
	return nil
}

func list(ctx context.Context, client todo.TasksClient) error {
	l, err := client.List(ctx, &todo.Void{})
	if err != nil {
		return fmt.Errorf("could not fetch tasks: %v", err)
	}

	fmt.Println("\nHere are your todos:")

	for _, t := range l.Tasks {
		if t.Done {
			fmt.Printf("[X]: ")
		} else {
			fmt.Printf("[ ]: ")
		}
		fmt.Println(t.Title)
	}
	return nil
}

func done(ctx context.Context, client todo.TasksClient, text string) error {
	t, err := client.Done(ctx, &todo.TaskTitle{Title: text})
	if err != nil {
		return fmt.Errorf("could not complete task: %v", err)
	}
	fmt.Println("completed task:", t.Title)
	return nil
}
