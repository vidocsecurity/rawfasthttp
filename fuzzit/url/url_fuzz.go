//go:build gofuzz
// +build gofuzz

package fuzz

import (
	"bytes"
)

func Fuzz(data []byte) int {
	u := rawfasthttp.AcquireURI()
	defer rawfasthttp.ReleaseURI(u)

	u.UpdateBytes(data)

	w := bytes.Buffer{}
	if _, err := u.WriteTo(&w); err != nil {
		return 0
	}

	return 1
}
