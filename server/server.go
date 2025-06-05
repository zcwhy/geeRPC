package server

import (
	"fmt"
	"geerpc/codec"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

type Server struct {
	serviceMap map[string]*serviceType
}

func NewServer() *Server {
	return &Server{
		serviceMap: make(map[string]*serviceType),
	}
}

func (s *Server) Run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicf("rpc server: listen address %v error: %v", address, err)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Println("rpc server: accept error:", err)
			return err
		}
		go s.serveConn(conn)
	}
}

func (s *Server) serveConn(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	f := codec.NewCodecFuncMap[codec.JsonType]
	s.serveRequests(f(conn))
}

func (s *Server) serveRequests(c codec.Codec) {
	wg := &sync.WaitGroup{}

	for {
		req := &codec.Header{}
		if err := c.ReadHeader(req); err != nil {
			log.Printf("rpc server: read rpc struct error: %v", err)
			break
		}

		service, method, err := s.findHandler(req.ServiceMethod)
		if err != nil {
			log.Printf("rpc server: find handler error: %v", err)
			return
		}

		argvType := method.arg
		argv := reflect.New(argvType)
		if err := c.ReadBody(argv.Interface()); err != nil {
			log.Println("rpc server: read argv err:", err)
		}

		wg.Add(1)
		go s.handleRequest(req, service, method, argv, c, wg)
	}
	wg.Wait()
}

func (s *Server) handleRequest(reqHeader *codec.Header, sType *serviceType, mType *MethodType, argv reflect.Value, c codec.Codec, wg *sync.WaitGroup) {
	defer wg.Done()

	reply := reflect.New(mType.reply.Elem())
	f := mType.method.Func
	f.Call([]reflect.Value{sType.rcvr, argv.Elem(), reply})

	if err := c.Write(reqHeader, reply.Interface()); err != nil {
		log.Printf("rpc server: write rpc struct error: %v", err)
	}
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

func (s *Server) findHandler(serviceMethod string) (*serviceType, *MethodType, error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		return nil, nil, fmt.Errorf("rpc server: service/method request ill-formed: %v", serviceMethod)
	}

	service, method := serviceMethod[:dot], serviceMethod[dot+1:]

	srvType, ok := s.serviceMap[service]
	if !ok {
		return nil, nil, fmt.Errorf("[findHandler] failed to find handler for serviceMethod %s", serviceMethod)
	}

	mType, ok := srvType.method[method]
	if !ok {
		return nil, nil, fmt.Errorf("[findHandler] failed to find handler for serviceMethod %s", serviceMethod)
	}

	return srvType, mType, nil
}
