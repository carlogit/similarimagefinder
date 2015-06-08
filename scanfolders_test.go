package main

import (
	"log"
	"path/filepath"
	"testing"
)

func TestGetJPGFilePaths(t *testing.T) {
	stop := make(chan struct{})
	defer close(stop)

	dir, err := filepath.Abs("testdata")
	if err != nil {
		log.Fatal(err)
	}

	pathsChan := getJPGFilePaths(stop, dir)

	index := 0
	paths := make([]string, 2)
	for path := range pathsChan {
		paths[index] = path
		index++
	}

	filename := paths[0]
	expectedPath, _ := filepath.Abs("testdata/soccerball.jpg")
	if filename != expectedPath {
		t.Errorf("Jpeg found is %s, want => %s", filename, expectedPath)
	}

	filename = paths[1]
	expectedPath, _ = filepath.Abs("testdata/subfolder/soccerball.jpeg")
	if filename != expectedPath {
		t.Errorf("Jpeg found is %s, want => %s", filename, expectedPath)
	}

}
