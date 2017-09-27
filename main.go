package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

// Filter declare message filter
type Filter func(m events.Message) bool

// Select stream elements by filter
func Select(in *<-chan events.Message, isSelected Filter, out *chan events.Message) {
	go func() {
		fmt.Println("goroutine started")
		for {
			select {
			case m := <-*in:
				if isSelected(m) {
					*out <- m
				}
			}
		}
	}()
}

func print(out *chan events.Message) {
	fmt.Println("print started")
	for {
		select {
		case m := <-*out:
			fmt.Print("print => ")
			fmt.Println(m.Status)
		}
	}
}

func dockerStream() <-chan events.Message {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
	messages, _ := cli.Events(context.Background(), types.EventsOptions{})
	return messages
}

func main() {
	var stream = dockerStream()
	execStream := make(chan events.Message)
	var execStart = func(m events.Message) bool {
		return strings.HasPrefix(m.Status, "exec_start")
	}
	Select(&stream, execStart, &execStream)
	print(&execStream)
}
