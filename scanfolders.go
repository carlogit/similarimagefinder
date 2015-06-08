package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func getJPGFilePaths(done <-chan struct{}, folder string) <-chan string {
	paths := make(chan string)
	go func() {
		defer close(paths)

		filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
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
