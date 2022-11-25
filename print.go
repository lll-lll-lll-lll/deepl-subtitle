package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrintlnJson(elements []*Element) {
	for _, e := range elements {
		var out bytes.Buffer
		b, _ := json.Marshal(e)
		err := json.Indent(&out, b, "", "  ")
		if err != nil {
			panic(err)

		}
		fmt.Println(out.String())
	}
}
