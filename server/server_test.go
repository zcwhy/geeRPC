package server

import (
	"reflect"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer()
	s.Register(TestStruct{})
	t.Log(s.findHandler("TestStruct.TestFunc1"))
}

func TestNewArgv(t *testing.T) {
	i := 1
	argvType := reflect.TypeOf(i)

	argv := reflect.New(argvType)
	t.Log(argv.Kind() == reflect.Ptr)
	t.Log(argv.Elem().Kind())
}

func TestReflect(t *testing.T) {
	//temp := "Hello world!"
	//reply := &temp
}
