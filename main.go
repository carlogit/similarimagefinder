package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/carlogit/phash"
)

func calculateHashes(done <-chan struct{}, root string) (<-chan *result, <-chan error) {
	// For each regular file, start a goroutine that calculates the file hashes and sends
	// the result on c.  Send the result of the walk on errc.
	c := make(chan *result)
	errc := make(chan error, 1)
	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println("cannot scan path %s, error: %d", path, err.Error())
				return nil
			}
			if !info.Mode().IsRegular() || filepath.Ext(path) != "jpg" || filepath.Ext(path) != "jpg"){
				return nil
			}
			wg.Add(1)
			go func() {
				result, err := buildResult(path)
				if err != nil {
					log.Println("cannot calculate hash for image %s, error: %d", path, err.Error())
				} else {
					select {
					case c <- result:
					case <-done:
					}
				}
				wg.Done()
			}()
			// Abort the walk if done is closed.
			select {
			case <-done:
				return errors.New("walk canceled")
			default:
				return nil
			}
		})
		// Walk has returned, so all calls to wg.Add are done.  Start a
		// goroutine to close c once all the sends are done.
		go func() {
			wg.Wait()
			close(c)
		}()
		// No select needed here, since errc is buffered.
		errc <- err
	}()
	return c, errc
}

func calculateAllHashes(root string) (map[string]*result, error) {
	// MD5All closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	c, errc := calculateHashes(done, root)

	m := make(map[string]*result)
	for r := range c {
		m[r.path] = r
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

type similarResult struct {
	path    string
	same    []string
	similar []string
}

func main() {
	//	m, err := MD5All(os.Args[1])
	m, err := calculateAllHashes("/home/carlo/Downloads/Camera/")
	if err != nil {
		fmt.Println(err)
		return
	}

	s := make(map[string]similarResult)
	for key, value := range m {
		same := make([]string, 0)
		similar := make([]string, 0)
		for key1, value1 := range m {
			if key == key1 {
				continue
			}

			if value.sha1 == value1.sha1 {
				same = append(same, key1)
			} else if phash.GetDistance(value.phash, value1.phash) <= 5 {
				similar = append(similar, key1)
			}
		}

		s[key] = similarResult{key, same, similar}
	}

	for key, value := range s {
		fmt.Printf("image: %s\n", key)
		fmt.Printf("  same:\n")
		for _, value1 := range value.same {
			fmt.Printf("    %s\n", value1)
		}

		fmt.Printf("  similar:\n")
		for _, value2 := range value.similar {
			fmt.Printf("    %s\n", value2)
		}
	}
}
