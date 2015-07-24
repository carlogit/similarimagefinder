package main

import (
	"testing"

	"github.com/carlogit/similarimagefinder/fingerprint"
	"log"
)

func TestBuildSimilarImagesList(t *testing.T) {
	mapPathHashes := make(map[string]string)
	mapPathHashes["pathA"] = "11111111"
	mapPathHashes["pathB"] = "01001111"
	mapPathHashes["pathC"] = "11111100"
	mapPathHashes["pathD"] = "10001111"
	mapPathHashes["pathE"] = "11111111"
	mapPathHashes["pathF"] = "00000000"

	similarImagesList := fingerprint.BuildSimilarImagesList(2, mapPathHashes)

	log.Println(similarImagesList)
	if len(similarImagesList) != 2 {
		t.Errorf("Number of dupsets is %d, want => %d", len(similarImagesList), 2)
	}

	if len(similarImagesList[0]) == 3 && len(similarImagesList[1]) == 2 {
		checkElements(t, similarImagesList[0], "pathA", "pathC", "pathE")
		checkElements(t, similarImagesList[1], "pathB", "pathD")
	} else if len(similarImagesList[0]) == 2 && len(similarImagesList[1]) == 3 {
		checkElements(t, similarImagesList[0], "pathB", "pathD")
		checkElements(t, similarImagesList[1], "pathA", "pathC", "pathE")
	} else {
		t.Errorf("Unexpected number of items in one or both of the expected dupsets")
	}
}

func checkElements(t *testing.T, dupSet map[string]bool, paths... string) {
	for _, path := range paths {
		if !dupSet[path] {
			t.Errorf("%s has not been found in dupset, values found %v, ", path, dupSet)
		}
	}
}
