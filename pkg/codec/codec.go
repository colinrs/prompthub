package codec

import (
	"encoding/json"

	"github.com/bytedance/sonic"
)

type Codec interface {
	Marshal(value interface{}) ([]byte, error)
	Unmarshal(b []byte, dst interface{}) error
}

var _ Codec = (*defaultCodec)(nil)

func NewCodec() Codec {
	return &defaultCodec{}
}

type defaultCodec struct{}

func (c *defaultCodec) Marshal(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

func (c *defaultCodec) Unmarshal(b []byte, dst interface{}) error {
	return json.Unmarshal(b, dst)
}

type sonicCodec struct{}

func (c *sonicCodec) Marshal(value interface{}) ([]byte, error) {
	return sonic.Marshal(value)
}

func (c *sonicCodec) Unmarshal(b []byte, dst interface{}) error {
	return sonic.Unmarshal(b, dst)
}

func NewSonicCodec() Codec {
	return &sonicCodec{}
}
