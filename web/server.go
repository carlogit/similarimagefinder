package web

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"os"
)

func StartWebService(port int) {
	log.Println(fmt.Sprintf("Starting service on port %d", port))
	http.HandleFunc("/delete", deleteHandler)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)

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

