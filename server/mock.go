package server

import "log"

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func MockServer() {
	s := NewServer()

	var foo Foo
	if err := s.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}

	s.Run(":8080")
}
