package codec

import "io"

type Codec interface {
	ReadHeader(*Header) error
	ReadBody(any) error
	Write(*Header, any) error
}

type NewCodecFunc func(conn io.ReadWriteCloser) Codec

type CodeType uint8

const (
	GobType  CodeType = 1 // gob
	JsonType CodeType = 2
)

var NewCodecFuncMap map[CodeType]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[CodeType]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
	NewCodecFuncMap[JsonType] = NewJsonCodec
}

type Header struct {
	ServiceMethod string
	Seq           uint64
	Error         error
}
