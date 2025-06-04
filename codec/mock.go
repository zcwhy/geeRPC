package codec

import (
	"bytes"
)

type MockConn struct {
	n   int
	buf *bytes.Buffer
}

func NewMockConn() *MockConn {
	return &MockConn{
		buf: bytes.NewBuffer(make([]byte, 0, 4096)),
	}
}

func (m MockConn) Read(p []byte) (n int, err error) {

	return m.buf.Read(p)
}

func (m MockConn) Write(p []byte) (n int, err error) {
	n, err = m.buf.Write(p)
	//fmt.Println(m.buf.Available())
	return n, err
}

func (m MockConn) Close() error {
	return nil
}
