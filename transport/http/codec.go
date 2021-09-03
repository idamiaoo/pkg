package http

import (
	"net/url"

	"github.com/go-playground/form/v4"
	jsoniter "github.com/json-iterator/go"
)

type FormCodec interface {
	Marshal(v interface{}) (values url.Values, err error)
	Unmarshal(values url.Values, v interface{}) error
}

type defaultFormCodec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

func NewDefaultFormCodec() FormCodec {
	encoder := form.NewEncoder()
	decoder := form.NewDecoder()
	encoder.SetTagName("json")
	decoder.SetTagName("json")
	return &defaultFormCodec{
		encoder: encoder,
		decoder: decoder,
	}
}

func (c defaultFormCodec) Marshal(v interface{}) (values url.Values, err error) {
	vs, err := c.encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return vs, nil
}

func (c defaultFormCodec) Unmarshal(values url.Values, v interface{}) error {
	if err := c.decoder.Decode(v, values); err != nil {
		return err
	}
	return nil
}

type JsonCodec interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

type defaultJsonCodec struct {
	jsoniter.API
}

func NewDefaultJsonCodec() JsonCodec {
	return &defaultJsonCodec{
		API: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}
