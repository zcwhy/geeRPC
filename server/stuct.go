package server

import "fmt"

type TestStruct struct{}

func (s TestStruct) TestFunc1(arg int, reply *string) error {
	*reply = fmt.Sprintf("Receive %d success!", arg)
	return nil
}

func (s TestStruct) TestFunc2(arg int) error {
	return nil
}

func (s TestStruct) TestFunc3(arg int, reply *string) (int, error) {
	return 0, nil
}

// second arg not a pointer
func (s TestStruct) TestFunc4(arg int, reply string) error {
	return nil
}
