package server

import (
	"fmt"
	"geerpc/codec"
	"log"
	"net"
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

		//req.ServiceMethod

		// TODO: now we don't know the type of request argv
		// day 1, just suppose it's string
		var argv string
		if err := c.ReadBody(&argv); err != nil {
			log.Println("rpc server: read argv err:", err)
		}

		fmt.Println("[server] 123", argv)

		fmt.Println("[server] handle request:", req)
		wg.Add(1)
		go s.handleRequest(req, c, wg)
	}
	wg.Wait()
}

func (s *Server) handleRequest(reqHeader *codec.Header, c codec.Codec, wg *sync.WaitGroup) {
	defer wg.Done()

	//cmd := reqHeader.ServiceMethod
	//strings.Split(cmd, ".")
	//
	//sType, mType, err := s.findHandler(reqHeader.ServiceMethod)
	//if err != nil {
	//	log.Printf("[handleRequest] error: %v", err)
	//	return
	//}
	//
	//arg := mType.newArg()
	//reply := mType.newReply()
	//
	//c.ReadBody(arg.Interface())

	//f := mType.method.Func
	//if mType.arg.Kind() != reflect.Ptr {
	//	arg = arg.Elem()
	//}
	//f.Call([]reflect.Value{sType.rcvr, arg, reply})

	if err := c.Write(reqHeader, "Hello, world"); err != nil {
		log.Printf("rpc server: write rpc struct error: %v", err)
	}
}

func (s *Server) findHandler(serviceMethod string) (*serviceType, *MethodType, error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		log.Println("")
	}

	service, method := serviceMethod[:dot], serviceMethod[dot+1:]

	srvType, ok := s.serviceMap[service]
	if !ok {
		return nil, nil, fmt.Errorf("[findHandler] failed to find handler for serviceMethod %s", serviceMethod)
	}

	mType, ok := srvType.method[method]
	if !ok {

	}

	return srvType, mType, nil
}
