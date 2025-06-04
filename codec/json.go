package codec

import (
	"bufio"
	"encoding/json"
	"io"
)

type jsonCodec struct {
	conn    io.ReadWriteCloser
	buf     *bufio.Writer
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &jsonCodec{
		conn:    conn,
		buf:     buf,
		encoder: json.NewEncoder(buf),
		decoder: json.NewDecoder(conn),
	}
}

func (c *jsonCodec) ReadHeader(h *Header) error {
	return c.decoder.Decode(h)
}

func (c *jsonCodec) ReadBody(a any) error {
	return c.decoder.Decode(a)
}

func (c *jsonCodec) Write(header *Header, body any) error {
	if err := c.encoder.Encode(header); err != nil {
		return err
	}

	if err := c.encoder.Encode(body); err != nil {
		return err
	}

	return c.buf.Flush()
}
