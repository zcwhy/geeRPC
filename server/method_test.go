package server

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type DummyService struct{}

func (s *DummyService) DummyFunc1(arg *int, reply int) error {
	return nil
}

func TestService_Register(t *testing.T) {
	s := NewServer()

	err := s.Register(&DummyService{})
	assert.NoError(t, err)
	//srvType, mType, err := s.findHandler("DummyService.DummyFunc1")
	//assert.NoError(t, err)
	//
	//assert.Equal(t, srvType.name, "DummyService")
	//assert.Equal(t, mType.name, "DummyFunc1")
}

func TestFilterMethods(t *testing.T) {
	methods := filterMethods(reflect.TypeOf(&DummyService{}))
	assert.Len(t, methods, 1)
	assert.NotNil(t, methods["DummyFunc1"])
	t.Log(methods["DummyFunc1"])
}

func Test(t *testing.T) {
	a := 3
	i := &a
	*i = 4

	t.Log(a)
}
