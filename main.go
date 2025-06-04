package main

import (
	"fmt"
	"geerpc/client"
	"geerpc/server"
	"time"
)

func main() {
	go func() {
		srv := server.NewServer()
		srv.Run(":8080")
	}()

	time.Sleep(1 * time.Second)
	c, err := client.Dial(":8080")
	if err != nil {
		panic(err)
	}

	var resp string
	//args := &server.Args{Num1: 1, Num2: 2}
	//time.Sleep(1 * time.Second)

	fmt.Println(c.Call("Foo.Sum", "123", &resp))
	fmt.Println(resp)
}
