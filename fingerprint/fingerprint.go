package fingerprint

import (
	"bytes"
	"io/ioutil"

	"github.com/carlogit/phash"
	"log"
)

func CalculatePhash(filepath string) (string, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	reader := bytes.NewReader(data)
	phash, err := phash.GetHash(reader)
	if err != nil {
		return "", err
	}

	return phash, nil
}

func BuildSimilarImagesList(distanceThreshold int, pathHashMap map[string]string) []map[string]bool {
	log.Println("Finding duplicate images...")
	similarImagesList := make([]map[string]bool, 0)
	processedImages := make(map[string]bool)
	for path, hash := range pathHashMap {
		if processedImages[path] {
			continue
		}
		log.Println("Processing: " + path + " - " + hash)
		processedImages[path] = true

		duplicates := make(map[string]bool)
		for otherPath, otherHash := range pathHashMap {
			if processedImages[otherPath] {
				continue
			}
			if phash.GetDistance(hash, otherHash) <= distanceThreshold {
				log.Println("duplicates found")
				duplicates[otherPath] = true
			}
		}
		if len(duplicates) > 0 {
			for k,v := range duplicates {
				processedImages[k] = v
			}

			duplicates[path] = true
			similarImagesList = append(similarImagesList, duplicates)
			log.Printf("Duplicate set found: %v\n", duplicates)
		}
	}

	return similarImagesList
}
