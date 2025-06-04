package codec

import (
	"net"
	"testing"
)

func TestCodec(t *testing.T) {
	s := RPCRequest{
		ServiceMethod: "foo.bar",
		Seq:           uint64(1),
	}

	mockConn := NewMockConn()

	gobCodecFn := NewCodecFuncMap[GobType]
	codec := gobCodecFn(mockConn)

	if err := codec.WriteRequest(&s); err != nil {
		t.Fatal(err)
	}

	var es = &RPCRequest{}
	if err := codec.ReadRequest(es); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", es)
}

func TestPipe(t *testing.T) {
	// 使用 net.Pipe 创建一个模拟的 TCP 连接对
	conn1, conn2 := net.Pipe()

	gobCodecFn := NewCodecFuncMap[GobType]

	s := RPCRequest{
		ServiceMethod: "foo.bar",
		Seq:           uint64(1),
	}

	go func() {
		codec := gobCodecFn(conn1)
		codec.WriteRequest(&s)
	}()

	var es = &RPCRequest{}
	codec := gobCodecFn(conn2)
	codec.ReadRequest(es)

	t.Logf("%+v\n", es)
}

func TestReadRequest(t *testing.T) {
	req := &RPCRequest{
		ServiceMethod: "foo.bar",
	}

	gobCodecFn := NewCodecFuncMap[GobType]

	conn1, conn2 := net.Pipe()
	codec := gobCodecFn(conn1)
	codec.WriteRequest(req)

	ss := &RPCRequest{}
	codec = gobCodecFn(conn2)
	codec.ReadRequest(ss)
	t.Logf("%+v\n", ss)
}
