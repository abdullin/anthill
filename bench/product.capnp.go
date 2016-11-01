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
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 4})
	return Product{st}, err
}

func NewRootProduct(s *capnp.Segment) (Product, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 4})
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

func (s Product) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s Product) SetId(v uint64) {
	s.Struct.SetUint64(0, v)
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
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 4}, sz)
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

const schema_85d3acc39d94e0f8 = "x\xda\\\x911K\xebP\x1c\xc5\xcf\xb97}y\x85" +
	"\xf2\x9aK\x03\xef\xf1\x96\x828T\xa1\xa2\xe8 ]*" +
	"\xf8\x05z\x07\x07\x07\x87\x98T\x08j\x12\x9a\x16E\x14" +
	"\x94Z\xa8\xe0\xa0X7\xfd\x06N.\x0e\xba\xf9\x11\xdc" +
	"\xc5Mps\xd3A\"QI\xd1\xe9r\x0f\xff\xff\xef" +
	"\x1c\xce\xdf\xba\x99\x13S\xb9\xbf\x02\xd0\xffr\xbf\x92z" +
	"\xb7V\xb7\xac\xca)\xf4\x1f2yy\x18\x9c\xdf^\xdc" +
	"\xf5\x90\xa3\x09\xa8\xe7'\xf5\x96\xbe\xaf\x1b`r}\x7f" +
	"2{U\xe9>\xfe\x9c4L\xa0\xb4\xc0\xcb\xd2R\xba" +
	"3\xbd\xc82QM\xa2V\xe8u\xdc\xf6\x84p\x9d(" +
	"\x88j\xf3kN\x1c\xfb+\xbe\xeb\x14\xdb~\x184H" +
	"\xfd[\x1a\x80A@\x8d\xfd\x07\xf4\xa8\xa4\x9e\x14$m" +
	"\xa6Zu\x1c\xd0\x15I=#(}\x8fy\x08\xe6\xc1" +
	"b\xe0\xac7Y\x80`\x01\xccL\xf8i\xd2h\x95?" +
	"\xfe)\xdd\xce\xe8;)}SR\xef\x0f\xe9{)}" +
	"[R\xf7\x05\x95\xa0M\x01\xa8\xde\x08\xa0w%\xf5\xa1" +
	"\xa0\x92\xc2\xa6\x04\xd4\xc12\xa0\xfb\x92z \xa8\x0ci" +
	"\xd3\x00\xd4\xf1\x16\xa0\x8f$\xf5\xd9\xf7pn\xe8e\xe1" +
	"\xccx\xb5\x93\x05\xf5\x9a\xb1\xdb\xf2\xa36L?\x0c2" +
	"\xd5\xfd*\x05u\xd7I[\xa15<\x07H\x0b|\x0f" +
	"\x00\x00\xff\xff\xea\xb9^`"

func init() {
	schemas.Register(schema_85d3acc39d94e0f8,
		0x952810103f3a833f,
		0xe78328b73893debb)
}
