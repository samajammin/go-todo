# go-todo

Simple TODO cli tool using gRPC, inspired by [justforfunc](https://youtu.be/_jQ3i_fyqGA)'s tutorial

## Installation

1.  Clone the repo within your GOPATH
2.  From the project root, run `go install ./cmd/todo`
3.  From the project root, create a file named tododb.pb to store your todos
4.  Add some todos! Get some shit done!

## Commands

### Add

Add a todo with `todo add [task]`

### List

List your todos with `todo list`

### Done

Mark a task as complete with `todo done [task]`
