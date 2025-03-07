// Package urlencode provides a urlencode codec
package urlencode

import (
	"encoding/json"

	pb "go.unistack.org/micro-proto/v4/codec"
	"go.unistack.org/micro/v4/codec"
	rutil "go.unistack.org/micro/v4/util/reflect"
	"google.golang.org/protobuf/types/known/structpb"
)

type urlencodeCodec struct {
	opts codec.Options
}

var _ codec.Codec = &urlencodeCodec{}

func (c *urlencodeCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	}

	uv, err := rutil.StructURLValues(v, "", []string{"protobuf", "json", "xml", "yaml"})
	if err != nil {
		return nil, err
	}

	return []byte(uv.Encode()), nil
}

func (c *urlencodeCodec) Unmarshal(b []byte, v interface{}, opts ...codec.Option) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = b
		return nil
	case *pb.Frame:
		m.Data = b
		return nil
	case *codec.RawMessage:
		*m = append((*m)[0:0], b...)
		return nil
	case codec.RawMessage:
		copy(m, b)
		return nil
	}

	mp, err := rutil.URLMap(string(b))
	if err != nil {
		return err
	}

	switch t := v.(type) {
	case *structpb.Value:

		buf, err := json.Marshal(mp)
		if err == nil {
			err = t.UnmarshalJSON(buf)
		}
		return err
	}

	return rutil.Merge(v, rutil.FlattenMap(mp), rutil.Tags([]string{"protobuf", "json", "xml", "yaml"}), rutil.SliceAppend(true))
}

func (c *urlencodeCodec) String() string {
	return "urlencode"
}

func NewCodec(opts ...codec.Option) *urlencodeCodec {
	return &urlencodeCodec{opts: codec.NewOptions(opts...)}
}
