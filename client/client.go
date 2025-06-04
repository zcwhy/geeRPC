package client

import (
	"fmt"
	"geerpc/codec"
	"log"
	"net"
)

var NextSeq uint64 = 0

type Client struct {
	conn    net.Conn
	cc      codec.Codec
	pending map[uint64]*Call
}

type Call struct {
	id            uint64
	serviceMethod string
	args          any
	reply         any
	done          chan *Call
	error         error
}

func Dial(addr string) (client *Client, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	return NewClient(conn)
}

func NewClient(conn net.Conn) (*Client, error) {
	c := &Client{
		conn:    conn,
		cc:      codec.NewJsonCodec(conn),
		pending: make(map[uint64]*Call),
	}
	// json send code type
	go c.receive()

	return c, nil
}

func (c *Client) Call(serviceMethod string, args any, reply any) error {
	cc := <-c.Go(serviceMethod, args, reply).done
	return cc.error
}

func (c *Client) Go(serviceMethod string, args any, reply any) *Call {
	call := &Call{
		id:            getSeq(),
		serviceMethod: serviceMethod,
		args:          args,
		reply:         reply,
		done:          make(chan *Call, 1),
	}

	c.send(call)

	return call
}

func (c *Client) send(call *Call) {
	req := codec.Header{
		Seq:           call.id,
		ServiceMethod: call.serviceMethod,
	}

	fmt.Println("[send] 111111", req, call.args)

	c.pending[call.id] = call
	if err := c.cc.Write(&req, call.args); err != nil {
		log.Println("rpc client:", err)
		call.error = err
	}
}

func (c *Client) receive() {
	for {
		respHeader := codec.Header{}
		if err := c.cc.ReadHeader(&respHeader); err != nil {
			log.Println("rpc client: read response error:", err)
			return
		}

		fmt.Println(respHeader)
		callId := respHeader.Seq
		call, ok := c.pending[callId]
		if !ok {

		}
		//call.error = respHeader.Error

		if err := c.cc.ReadBody(&call.reply); err != nil {
			log.Println("rpc client: read response error:", err)
			return
		}
		delete(c.pending, callId)
		call.done <- call
	}
}

func getSeq() uint64 {
	NextSeq++
	return NextSeq
}
