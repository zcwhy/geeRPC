package main

import (
	"context"
	"fmt"
	"geerpc/client"
	"geerpc/server"
	"time"
)

func main() {
	go func() {
		srv := server.NewServer()
		srv.Register(server.TestStruct{})
		srv.Run(":8080")
	}()

	time.Sleep(1 * time.Second)
	c, err := client.Dial(":8080")
	if err != nil {
		panic(err)
	}

	args := 1
	var resp string

	fmt.Println(c.Call(context.Background(), "TestStruct.TestFunc1", args, &resp))
	fmt.Println(resp)
}
