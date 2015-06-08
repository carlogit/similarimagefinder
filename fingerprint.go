package main

import (
	"bytes"
	"io/ioutil"

	"github.com/carlogit/phash"
)

type pathHashes struct {
	path  string
	phash string
}

func buildResult(filepath string) (*pathHashes, error) {
	data, err := ioutil.ReadFile(filepath)
	reader := bytes.NewReader(data)

	phashData, err := phash.GetHash(reader)
	if err != nil {
		return nil, err
	}

	return &pathHashes{filepath, phashData}, nil
}
