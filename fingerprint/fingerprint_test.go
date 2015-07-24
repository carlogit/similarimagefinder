package fingerprint

import (
	"testing"
)

func TestBuildResult(t *testing.T) {
	result, err := CalculatePhash("../testdata/unexistantFile.jpg")

	if err == nil {
		t.Errorf("Error expected due to file does not exist")
	}

	result, err = CalculatePhash("../testdata/soccerball.jpg")

	if err != nil {
		t.Errorf("No error expected, error found %s", err.Error())
	} else {
		if len(result) != 64 {
			t.Errorf("A phash value of %d is expected, found %d", 64, len(result))
		}
	}
}
