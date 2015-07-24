package folderscanner

import (
	"log"
	"os"
	"path/filepath"
)

func GetJPGFilePaths(folder string) <-chan string {
	pathsChan := make(chan string)
	go func() {
		defer close(pathsChan)

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println("cannot scan path %s, error: %d", path, err.Error())
				return nil
			}
			if info.Mode().IsRegular() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg") {
				pathsChan <- path
			}

			return nil
		})
	}()

	return pathsChan
}
