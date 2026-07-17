package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/shibukawa/tinybind-go/jsonbind"
)

type Note struct {
	Text string `json:"text"`
	N    int    `json:"n"`
}

func main() {
	decoded, err := jsonbind.DecodeJSON[Note](strings.NewReader(`{"text":"tiny","n":7}`))
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	if err := jsonbind.EncodeJSON(&out, decoded); err != nil {
		panic(err)
	}
	fmt.Print(out.String())
}
