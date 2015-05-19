package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"

	"github.com/carlogit/phash"
)

type pathHashes struct {
	path  string
	sha1  string
	phash string
}

func calculateSha1(reader io.Reader) (string, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func buildResult(filepath string) (*pathHashes, error) {
	data, err := ioutil.ReadFile(filepath)
	reader := bytes.NewReader(data)

	sha1Data, err := calculateSha1(reader)
	if err != nil {
		return nil, err
	}

	reader.Seek(0, 0)

	phashData, err := phash.GetHash(reader)
	if err != nil {
		return nil, err
	}

	return &pathHashes{filepath, sha1Data, phashData}, nil
}
