package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/carlogit/similarimagefinder/web"
	"github.com/carlogit/similarimagefinder/fingerprint"
)

type imageResult struct {
	Path          string
	SimilarImages []string
}

func main() {
	threshold := flag.Int("threshold", 8, "hamming distance threshold.")
	folder := flag.String("folder", "", "full path for the folder to scan.")

	port := flag.Int("port", 9559, "port number for the service.")
	outputFile := flag.String("outFile", "similarimages.html", "full path for html output file.")

	flag.Parse()

	if *folder == "" {
		log.Fatalln("No full path value has been provided for argument folderToScan.")
		return
	}

	pathHashMap := fingerprint.CalculateHashes(*folder)

	similarImagesList := fingerprint.BuildSimilarImagesList(*threshold, pathHashMap)

	log.Println(similarImagesList);

	if len(similarImagesList) == 0 {
		log.Println("No duplicate images have been found.")
		return
	}

	err := createFileWithResults(similarImagesList, *outputFile, *port)
	if err != nil {
		return
	}

	log.Println("File has been created: " + *outputFile + ", please open file using a web browser.")

	web.StartWebService(*port)
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
