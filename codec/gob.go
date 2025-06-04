package codec

import (
	"bufio"
	"encoding/gob"
	"io"
)

type gobCodec struct {
	conn    io.ReadWriteCloser
	buf     *bufio.Writer
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &gobCodec{
		conn:    conn,
		buf:     buf,
		encoder: gob.NewEncoder(conn),
		decoder: gob.NewDecoder(conn),
	}
}

func (g *gobCodec) ReadHeader(rpcStruct *Header) error {
	//TODO implement me
	panic("implement me")
}

func (g *gobCodec) ReadBody(a any) error {
	//TODO implement me
	panic("implement me")
}

func (g *gobCodec) Write(rpcStruct *Header, a any) error {
	//TODO implement me
	panic("implement me")
}
