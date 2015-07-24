package fingerprint
import (
	"log"
	"sync"

	"github.com/carlogit/similarimagefinder/folderscanner"
)

const numberOfGoRoutines = 15

type filePathHash struct {
	path  string
	phash string
}

func calculatePhashForFilePath(pathsChan <-chan string, c chan <- filePathHash) {
	for filePath := range pathsChan {
		phash, err := CalculatePhash(filePath)
		if err != nil {
			log.Println("Cannot calculate Phash for file: " + filePath)
		} else {
			c <- filePathHash{filePath, phash}
		}
	}
}

func CalculateHashes(root string) map[string]string {
	log.Println("Calculating phash fingerprint for images...")

	pathsChan := folderscanner.GetJPGFilePaths(root)

	pathHashChan := make(chan filePathHash)
	var wg sync.WaitGroup
	wg.Add(numberOfGoRoutines)
	for i := 0; i < numberOfGoRoutines; i++ {
		go func() {
			calculatePhashForFilePath(pathsChan, pathHashChan)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(pathHashChan)
	}()

	pathHashMap := make(map[string]string)
	for pathHash := range pathHashChan {
		pathHashMap[pathHash.path] = pathHash.phash
	}
	return pathHashMap
}


