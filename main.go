package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"

	"github.com/carlogit/phash"
)

const numDigesters = 20

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
	log.Println("Calculating fingerprint for images...")
	done := make(chan struct{})
	defer close(done)

	paths := getJPGFilePaths(done, root)

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

type imageResult struct {
	Path          string
	SimilarImages []string
}

func main() {

	threshold := flag.Int("threshold", 5, "hamming distance threshold.")
	folderToScan := flag.String("folderToScan", "", "full path for the folder to scan.")

	port := flag.Int("port", 9559, "port number for the service.")
	outputFile := flag.String("outFile", "similarimages.html", "full path for html output file.")

	flag.Parse()

	if *folderToScan == "" {
		log.Fatalln("No full path value has been provided for argument folderToScan.")
		return
	}

	mapPathHashes := calculateHashes(*folderToScan)

	similarImagesList := buildSimilarImagesList(*threshold, mapPathHashes)

	if len(similarImagesList) == 0 {
		log.Println("No duplicate images have been found.")
		return
	}

	err := createFileWithResults(similarImagesList, *outputFile, *port)
	if err != nil {
		return
	}

	log.Println("File has been created: " + *outputFile + ", please open file.")

	log.Println(fmt.Sprintf("Starting service on port %d", *port))
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), nil)

	//	for _, similarImage := range similarImagesMap {
	//		outputFile.WriteString(similarImage.Path + "\n")
	//		for _, similar := range similarImage.SimilarImages {
	//			outputFile.WriteString("  - " + similar + "\n")
	//		}
	//		outputFile.WriteString("\n")
	//	}

	//	outputFile.WriteString("\n")
	//	outputFile.WriteString("\n")
	//	outputFile.WriteString(string(b))

	//	outputFile.Sync()

	//	if err := json.NewEncoder(w).Encode(similarImagesMap); err != nil {
	//		log.Println(err)
	//		http.Error(w, "oops", http.StatusInternalServerError)
	//	}
}

func createFileWithResults(similarImagesList []map[string]bool, output string, port int) error {
	log.Println("Creating file (similarimages.html) with duplicate images...")
	results := buildDuplicateResults(similarImagesList)

	outputFile, err := os.Create(output)
	if err != nil {
		log.Fatalln("cannot create output file: %s. Error: %s", output, err.Error())
		return err
	}
	defer outputFile.Close()

	b, err := json.Marshal(results)
	if err != nil {
		log.Fatalln("cannot create JSON object from results", err.Error())
		return err
	}

	t, err := template.New("index.html").Delims("<<", ">>").ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalln("error parsing template file", err.Error())
		return err
	}

	templateData := struct {
		Images string
		Port   int
	}{string(b), port}

	err = t.Execute(outputFile, templateData)
	if err != nil {
		log.Fatalln("error processing template file", err.Error())
		return err
	}

	return nil
}

func buildDuplicateResults(similarImagesList []map[string]bool) [][]string {
	results := make([][]string, len(similarImagesList))
	for x := 0; x < len(similarImagesList); x++ {
		results[x] = make([]string, 0)
		for key, _ := range similarImagesList[x] {
			results[x] = append(results[x], key)
		}
	}

	return results
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	folderPath := r.URL.Query().Get("folderPath")
	callback := r.URL.Query().Get("callback")

	out := struct {
		Message string
		Status  bool
	}{}

	err := os.Remove(folderPath)
	if err != nil {
		out.Message = err.Error()
		out.Status = false
		log.Println(err.Error())
	} else {
		out.Message = "Image deleted."
		out.Status = true
	}

	jsonBytes, err := json.Marshal(out)
	if err != nil {
		log.Println(err)
		http.Error(w, "oops", http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s(%s)", callback, jsonBytes)
}

func buildSimilarImagesArray(distanceThreshold int, mapPathHashes map[string]*pathHashes) []imageResult {
	similarImagesArray := make([]imageResult, 0)
	for path, pathHashes := range mapPathHashes {
		similarImages := make([]string, 0)
		for otherPath, otherSimilarResult := range mapPathHashes {
			if path == otherPath {
				continue
			}

			if phash.GetDistance(pathHashes.phash, otherSimilarResult.phash) <= distanceThreshold {
				similarImages = append(similarImages, otherPath)
			}
		}

		if len(similarImages) > 0 {
			similarImagesArray = append(similarImagesArray, imageResult{path, similarImages})
		}
	}

	return similarImagesArray
}

func buildSimilarImagesList(distanceThreshold int, mapPathHashes map[string]*pathHashes) []map[string]bool {
	log.Println("Finding duplicate images...")
	similarImagesList := make([]map[string]bool, 0)
	for path, pathHashes := range mapPathHashes {
		duplicates := make(map[string]bool)
		for otherPath, otherSimilarResult := range mapPathHashes {
			if path == otherPath {
				continue
			}
			if phash.GetDistance(pathHashes.phash, otherSimilarResult.phash) <= distanceThreshold {
				if !duplicates[path] {
					duplicates[path] = true
				}
				if !duplicates[otherPath] {
					duplicates[otherPath] = true
				}
			}
		}
		if len(duplicates) > 0 {
			similarImagesList = addDuplicatesToMatches(duplicates, similarImagesList)
		}
	}

	return similarImagesList
}

func addDuplicatesToMatches(duplicates map[string]bool, similarImagesList []map[string]bool) []map[string]bool {
	compressed := false
	for _, matchSet := range similarImagesList {
		merge := false
		for key, _ := range duplicates {
			if matchSet[key] {
				merge = true
				break
			}
		}
		if merge {
			for key, _ := range duplicates {
				if matchSet[key] {
					matchSet[key] = true
				}
			}
			compressed = true
			break
		}
	}

	if !compressed {
		similarImagesList = append(similarImagesList, duplicates)
	}

	return similarImagesList
}
