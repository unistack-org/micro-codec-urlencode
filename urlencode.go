// Package urlencode provides a urlencode codec
package urlencode

import (
	"io"
	"io/ioutil"

	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type urlencodeCodec struct{}

func (c *urlencodeCodec) Marshal(b interface{}) ([]byte, error) {
	switch m := b.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	v, err := rutil.StructURLValues(b, "", []string{"protobuf", "json"})
	if err != nil {
		return nil, err
	}

	return []byte(v.Encode()), nil
}

func (c *urlencodeCodec) Unmarshal(b []byte, v interface{}) error {
	if len(b) == 0 {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = b
		return nil
	}

	mp, err := rutil.URLMap(string(b))
	if err != nil {
		return err
	}

	return rutil.Merge(v, rutil.FlattenMap(mp), rutil.Tags([]string{"protobuf", "json"}), rutil.SliceAppend(true))
}

func (c *urlencodeCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *urlencodeCodec) ReadBody(conn io.Reader, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		} else if len(buf) == 0 {
			return nil
		}
		m.Data = buf
		return nil
	}

	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}

	return c.Unmarshal(buf, b)
}

func (c *urlencodeCodec) Write(conn io.Writer, m *codec.Message, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	}

	buf, err := c.Marshal(b)
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)

	return err
}

func (c *urlencodeCodec) String() string {
	return "json"
}

func NewCodec() codec.Codec {
	return &urlencodeCodec{}
}
