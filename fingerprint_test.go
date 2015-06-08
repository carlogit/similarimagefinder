package main

import (
	"os"
	"testing"
)

func TestBuildResult(t *testing.T) {
	result, err := buildResult("testdata/unexistantFile.jpg")

	if err == nil {
		t.Errorf("Error expected due to file does not exist")
	}

	result, err = buildResult("testdata/soccerball.jpg")

	if err != nil {
		t.Errorf("No error expected, error found %s", err.Error())
	} else {
		expectedPath := "testdata/soccerball.jpg"
		if result.path != expectedPath {
			t.Errorf("Path is %s, want => %s", result.path, expectedPath)
		}

		if result.phash == "" {
			t.Errorf("A non empty phash value is expected")
		}
	}
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	return file
}
