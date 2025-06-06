package client

import (
	"context"
	"errors"
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

func Dial(addr string, opts ...OptionFunc) (client *Client, err error) {
	option := &Options{}
	for _, opt := range opts {
		opt(option)
	}

	conn, err := net.DialTimeout("tcp", addr, option.DialTimeOut)

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

func (c *Client) Call(ctx context.Context, serviceMethod string, args any, reply any) error {
	call := c.Go(serviceMethod, args, reply)

	select {
	case <-ctx.Done():
		return errors.New("rpc client: call failed: " + ctx.Err().Error())
	case <-call.done:
		return call.error
	}
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
