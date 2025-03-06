package securedtouch

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

const (
	key = "eG9yLWVuY3J5cHRpb24"
)

func Encrypt(data string) []byte {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, -1)
	if err != nil {
		fmt.Println(err)
	}
	w.Write([]byte(data))
	w.Close()
	zipped := b.Bytes()
	var result []byte
	for i := 0; i < len(zipped); i++ {
		result = append(result, byte(rune(zipped[i])^[]rune(key)[i%len(key)]))
	}
	fmt.Println(result)
	return result
}
