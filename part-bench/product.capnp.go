package main

// AUTO GENERATED - DO NOT EDIT

import (
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type Classification struct{ capnp.Struct }

// Classification_TypeID is the unique identifier for the type Classification.
const Classification_TypeID = 0x952810103f3a833f

func NewClassification(s *capnp.Segment) (Classification, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Classification{st}, err
}

func NewRootClassification(s *capnp.Segment) (Classification, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Classification{st}, err
}

func ReadRootClassification(msg *capnp.Message) (Classification, error) {
	root, err := msg.RootPtr()
	return Classification{root.Struct()}, err
}

func (s Classification) String() string {
	str, _ := text.Marshal(0x952810103f3a833f, s.Struct)
	return str
}

func (s Classification) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s Classification) SetId(v uint64) {
	s.Struct.SetUint64(0, v)
}

func (s Classification) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Classification) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Classification) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Classification) SetName(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

// Classification_List is a list of Classification.
type Classification_List struct{ capnp.List }

// NewClassification creates a new list of Classification.
func NewClassification_List(s *capnp.Segment, sz int32) (Classification_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Classification_List{l}, err
}

func (s Classification_List) At(i int) Classification { return Classification{s.List.Struct(i)} }

func (s Classification_List) Set(i int, v Classification) error { return s.List.SetStruct(i, v.Struct) }

// Classification_Promise is a wrapper for a Classification promised by a client call.
type Classification_Promise struct{ *capnp.Pipeline }

func (p Classification_Promise) Struct() (Classification, error) {
	s, err := p.Pipeline.Struct()
	return Classification{s}, err
}

type Product struct{ capnp.Struct }

// Product_TypeID is the unique identifier for the type Product.
const Product_TypeID = 0xe78328b73893debb

func NewProduct(s *capnp.Segment) (Product, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 4})
	return Product{st}, err
}

func NewRootProduct(s *capnp.Segment) (Product, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 4})
	return Product{st}, err
}

func ReadRootProduct(msg *capnp.Message) (Product, error) {
	root, err := msg.RootPtr()
	return Product{root.Struct()}, err
}

func (s Product) String() string {
	str, _ := text.Marshal(0xe78328b73893debb, s.Struct)
	return str
}

func (s Product) Code() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Product) HasCode() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Product) CodeBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Product) SetCode(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Product) Sku() (string, error) {
	p, err := s.Struct.Ptr(1)
	return p.Text(), err
}

func (s Product) HasSku() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Product) SkuBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	return p.TextBytes(), err
}

func (s Product) SetSku(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(1, t.List.ToPtr())
}

func (s Product) Description() (string, error) {
	p, err := s.Struct.Ptr(2)
	return p.Text(), err
}

func (s Product) HasDescription() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s Product) DescriptionBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(2)
	return p.TextBytes(), err
}

func (s Product) SetDescription(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(2, t.List.ToPtr())
}

func (s Product) Classification() (Classification, error) {
	p, err := s.Struct.Ptr(3)
	return Classification{Struct: p.Struct()}, err
}

func (s Product) HasClassification() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s Product) SetClassification(v Classification) error {
	return s.Struct.SetPtr(3, v.Struct.ToPtr())
}

// NewClassification sets the classification field to a newly
// allocated Classification struct, preferring placement in s's segment.
func (s Product) NewClassification() (Classification, error) {
	ss, err := NewClassification(s.Struct.Segment())
	if err != nil {
		return Classification{}, err
	}
	err = s.Struct.SetPtr(3, ss.Struct.ToPtr())
	return ss, err
}

// Product_List is a list of Product.
type Product_List struct{ capnp.List }

// NewProduct creates a new list of Product.
func NewProduct_List(s *capnp.Segment, sz int32) (Product_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 4}, sz)
	return Product_List{l}, err
}

func (s Product_List) At(i int) Product { return Product{s.List.Struct(i)} }

func (s Product_List) Set(i int, v Product) error { return s.List.SetStruct(i, v.Struct) }

// Product_Promise is a wrapper for a Product promised by a client call.
type Product_Promise struct{ *capnp.Pipeline }

func (p Product_Promise) Struct() (Product, error) {
	s, err := p.Pipeline.Struct()
	return Product{s}, err
}

func (p Product_Promise) Classification() Classification_Promise {
	return Classification_Promise{Pipeline: p.Pipeline.GetPipeline(3)}
}

const schema_85d3acc39d94e0f8 = "x\xdaL\x90\xbfJ\xc3`\x14\xc5\xcf\xf9\xbe\xd6X(" +
	"m>ZP\\\x04q\xa8BE\xc1A\xbaT\xf0\x05" +
	"z\x07\x07A\xc1\x98T\x1a\xacIlR\x04q+>" +
	"\x81:\xfa\x06\xba\xb888\xfa\x08\xee\xe2Vps\xd3" +
	"A\"\xf1Ot\xba\xdc\xc3\xbd\xbfs\xcf\xb5o\xd6\xd5" +
	"JqJ\x012]\x9cH\xdb\xa3V\xdb\xb6\x1b\x97\x90" +
	"\x0a\x99\xbe=_\\=\\?\x9e\xa1H\x0b0\xaf/" +
	"\xe6#\xab\xef\xc7`z\xfft\xbev\xd7\x18\x8da*" +
	"\xff\x07\x0b\x16P\xdb\xe4mm'[\xa9mq\x8cf" +
	"\x1a\x0dBo\xe8&K\xcau\xa2 jm\xf4\x9d8" +
	"\xf6\xf7}\xd7\xa9&~\x18tH\x99\xd4\x05\xa0@\xc0" +
	",\xcc\x002\xaf)\xcb\x8ad\x9d\x99\xd6\\\x04\xa4\xa1" +
	")\xab\x8a\xda\xf7X\x82b\x09\xac\x06\xcea\x97e(" +
	"\x96\xc1\xdc\x84\xdf&\x9d\xc1\xecW\x9f\xd1\xed\x9c\xeed" +
	"\xa4mM\xe9)\x9a_|w\x0e\x90]M\xe9+\x1a" +
	"\xa5\xeaT\x80\xf1\xf7\x00\xe9iJ\xa2h\xb4\xaeS\x03" +
	"\xe6\xe8\x04\x90HSN\x15\xabn\xe8\xe5\xeeV|0" +
	"\xcc/\xf1\xba\xb1;\xf0\xa3\x04\x96\x1f\x06\xb9\xea\xfe\xa4" +
	"F\xdbu\xb2\xd8\xb4\xff\xbe\x0d\xd2\x06?\x03\x00\x00\xff" +
	"\xfft\x99Y\xe1"

func init() {
	schemas.Register(schema_85d3acc39d94e0f8,
		0x952810103f3a833f,
		0xe78328b73893debb)
}
