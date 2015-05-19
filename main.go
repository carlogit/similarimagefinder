package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/carlogit/phash"
)

const numDigesters = 20
const threshold = 5

func walkFiles(done <-chan struct{}, root string) <-chan string {
	paths := make(chan string)
	go func() {
		defer close(paths)

		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println("cannot scan path %s, error: %d", path, err.Error())
				return nil
			}
			if !info.Mode().IsRegular() || (filepath.Ext(path) != ".jpg" && filepath.Ext(path) != ".jpeg") {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk canceled")
			}
			return nil
		})
	}()

	return paths
}

func digester(done <-chan struct{}, paths <-chan string, c chan<- *pathHashes) {
	for path := range paths {
		result, err := buildResult(path)
		if err == nil {
			select {
			case c <- result:
			case <-done:
				return
			}
		}
	}
}

func calculateHashes(root string) map[string]*pathHashes {
	done := make(chan struct{})
	defer close(done)

	paths := walkFiles(done, root)

	pathHashChan := make(chan *pathHashes)
	var wg sync.WaitGroup
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(done, paths, pathHashChan)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(pathHashChan)
	}()

	mapPathHash := make(map[string]*pathHashes)
	for pathHash := range pathHashChan {
		mapPathHash[pathHash.path] = pathHash
	}
	return mapPathHash
}

type similarResult struct {
	path    string
	same    []string
	similar []string
}

func main() {
	distanceThreshold := flag.Int("threshold", 5, "hamming distance to use for phash image similarity")
	folderToScan := flag.String("folder", "", "absolute path for the folder to scan for image similarity")

	flag.Parse()

	mapPathHashes := calculateHashes(*folderToScan)
	similarImagesMap := buildSimilarImagesMap(*distanceThreshold, mapPathHashes)

	outputFile, err := os.Create("similarimgage.txt")
	if err != nil {
		log.Fatalln("cannot create output file: %s. Error: %s", "similarimgage.txt", err.Error())
		return
	}
	defer outputFile.Close()

	for path, similarResult := range similarImagesMap {
		outputFile.WriteString(path + "\n")
		fmt.Printf("image: %s\n", path)
		fmt.Printf("  same:\n")
		for _, sameImage := range similarResult.same {
			fmt.Printf("    %s\n", sameImage)
			outputFile.WriteString("  * " + sameImage + "\n")
		}

		fmt.Printf("  similar:\n")
		for _, similarImage := range similarResult.similar {
			fmt.Printf("    %s\n", similarImage)
			outputFile.WriteString("  - " + similarImage + "\n")
		}

		outputFile.WriteString("\n")
	}

	outputFile.Sync()
}

func buildSimilarImagesMap(distanceThreshold int, mapPathHashes map[string]*pathHashes) map[string]similarResult {
	similarImagesMap := make(map[string]similarResult)
	for path, pathHashes := range mapPathHashes {
		sameImage := make([]string, 0)
		similarImage := make([]string, 0)
		for otherPath, otherSimilarResult := range mapPathHashes {
			if path == otherPath {
				continue
			}

			if pathHashes.sha1 == otherSimilarResult.sha1 {
				sameImage = append(sameImage, otherPath)
			} else if phash.GetDistance(pathHashes.phash, otherSimilarResult.phash) <= distanceThreshold {
				similarImage = append(similarImage, otherPath)
			}
		}

		similarImagesMap[path] = similarResult{path, sameImage, similarImage}
	}

	return similarImagesMap
}
