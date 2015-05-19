package main

import (
	"os"
	"testing"
)

func TestCalculateSha1(t *testing.T) {
	file1 := openFile("testdata/soccerball.jpg")
	defer file1.Close()

	sha1, err := calculateSha1(file1)
	if err != nil {
		t.Errorf("no error expected, error found %s", err.Error())
	}

	expectedSha1 := "349adc909b6bdf8acfec368a07593e95f5939fa0"
	if sha1 != expectedSha1 {
		t.Errorf("sha1 hash value %s, want %s", sha1, expectedSha1)
	}
}

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

		if result.sha1 == "" {
			t.Errorf("A non empty sha1 value is expected")
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
