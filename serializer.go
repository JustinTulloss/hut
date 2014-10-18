package hut

import (
	"encoding/json"
	"io"
	"net/http"
)

type Encoder interface {
	Encode(interface{}) error
}

type Decoder interface {
	Decode(interface{}) error
}

type Serializer interface {
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
	MimeType() string
}

type JsonSerializer struct{}

func (*JsonSerializer) NewEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func (*JsonSerializer) NewDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func (*JsonSerializer) MimeType() string {
	return "application/json; codec=utf-8"
}

var mimeSerializers = map[string]Serializer{
	"application/json": &JsonSerializer{},
}

func GetSerializer(r *http.Request) Serializer {
	return mimeSerializers["application/json"]
}
