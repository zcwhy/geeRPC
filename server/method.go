package server

import (
	"go/ast"
	"log"
	"reflect"
)

type serviceType struct {
	name   string
	rcvr   reflect.Value
	method map[string]*MethodType
}

type MethodType struct {
	name   string
	method reflect.Method
	arg    reflect.Type
	reply  reflect.Type
}

func (m *MethodType) newArg() reflect.Value {
	var v reflect.Value
	if m.arg.Kind() == reflect.Ptr {
		v = reflect.New(m.arg.Elem())
	} else {
		v = reflect.New(m.arg)
	}

	return v
}

// reply的类型一定是ptr
func (m *MethodType) newReply() reflect.Value {
	return reflect.New(m.reply.Elem())
}

func (s *Server) Register(rcvr any) error {
	t := reflect.TypeOf(rcvr)
	v := reflect.ValueOf(rcvr)

	service := &serviceType{
		rcvr:   reflect.ValueOf(rcvr),
		method: filterMethods(t),
	}
	service.name = t.Name()
	if v.Kind() == reflect.Ptr {
		service.name = reflect.Indirect(v).Type().Name()
	}
	s.serviceMap[service.name] = service
	return nil
}

func filterMethods(typ reflect.Type) map[string]*MethodType {
	var methods = make(map[string]*MethodType)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		mType := method.Type

		if !method.IsExported() {
			continue
		}

		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}

		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}

		methods[method.Name] = &MethodType{
			method: method,
			arg:    argType,
			reply:  replyType,
		}

		log.Printf("rpc server: register %s\n", method.Name)
	}

	return methods
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}
