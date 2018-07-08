package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/golang/protobuf/proto"

	"github.com/sbrichards/go-todo/todo"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type length int64

const (
	sizeOfLength = 8
	dbPath       = "tododb.pb"
)

var endianness = binary.LittleEndian

type taskServer struct{}

func main() {
	srv := grpc.NewServer()
	var tasks taskServer
	todo.RegisterTasksServer(srv, tasks)
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("could not connect to port: %v", err)
	}
	log.Fatal(srv.Serve(l))
}

func (s taskServer) Add(ctx context.Context, taskTitle *todo.TaskTitle) (*todo.Task, error) {
	if len(taskTitle.Title) < 1 {
		return nil, fmt.Errorf("cannot add empty task")
	}

	task := &todo.Task{
		Title: taskTitle.Title,
		Done:  false,
	}

	b, err := proto.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("could not encode task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %v", dbPath, err)
	}

	if err := binary.Write(f, endianness, length(len(b))); err != nil {
		return nil, fmt.Errorf("could not encode length of message: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, fmt.Errorf("could not write task to file: %v", err)
	}

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("could not close file %s: %v", dbPath, err)
	}

	return task, nil
}

func (s taskServer) List(ctx context.Context, void *todo.Void) (*todo.TaskList, error) {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %v", dbPath, err)
	}

	var taskList todo.TaskList
	for {
		if len(b) == 0 {
			return &taskList, nil
		} else if len(b) < sizeOfLength {
			return nil, fmt.Errorf("remaining odd %d bytes", len(b))
		}

		var l length
		if err := binary.Read(bytes.NewReader(b[:sizeOfLength]), endianness, &l); err != nil {
			return nil, fmt.Errorf("could not decode message length: %v", err)
		}
		b = b[sizeOfLength:]

		var task todo.Task
		if err := proto.Unmarshal(b[:l], &task); err != nil {
			return nil, fmt.Errorf("could not read task: %v", err)
		}
		b = b[l:]
		taskList.Tasks = append(taskList.Tasks, &task)
	}
}
