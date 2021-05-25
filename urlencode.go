// Package urlencode provides a urlencode codec
package urlencode

import (
	"io"

	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type urlencodeCodec struct{}

const (
	flattenTag = "flatten"
)

func (c *urlencodeCodec) Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	if m, ok := v.(*codec.Frame); ok {
		return m.Data, nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	uv, err := rutil.StructURLValues(v, "", []string{"protobuf", "json"})
	if err != nil {
		return nil, err
	}

	return []byte(uv.Encode()), nil
}

func (c *urlencodeCodec) Unmarshal(b []byte, v interface{}) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	if m, ok := v.(*codec.Frame); ok {
		m.Data = b
		return nil
	}

	mp, err := rutil.URLMap(string(b))
	if err != nil {
		return err
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	return rutil.Merge(v, rutil.FlattenMap(mp), rutil.Tags([]string{"protobuf", "json"}), rutil.SliceAppend(true))
}

func (c *urlencodeCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *urlencodeCodec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}

	return c.Unmarshal(buf, v)
}

func (c *urlencodeCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := c.Marshal(v)
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)
	return err
}

func (c *urlencodeCodec) String() string {
	return "urlencode"
}

func NewCodec() codec.Codec {
	return &urlencodeCodec{}
}
