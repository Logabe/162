//TODO: Fix up the index-out-of-bounds stuff

package main

import (
	"errors"
	"fmt"

	//"io"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"path/filepath"
)

type Data struct {
    Latest int `json:"latest"`
	Content []Entry `json:"content"`

}

type Entry struct {
	Location string `json:"location"`
	ContentType string `json:"content-type"`
	ContentDisposition string `json:"content-disposition"`
	ExecutionMethod string `json:"execution-method"`
	ExecutionData []string `json:"execution-data"`
}

var index Data

func main() {
	//Load the index file?
	data, err := os.ReadFile("index.json")
	if err != nil {
		log.Fatal(err)
	}

	//Unpack the JSON
	err = json.Unmarshal(data, &index)

	// '/' works as a default (use to acccept all requests unless another, more specific, handler has been set up.)
	http.HandleFunc("/", getPage)

	//NOTE: https://stackoverflow.com/questions/46992030/how-to-set-up-https-on-golang-web-server requies an Domain Name.
	err = http.ListenAndServe(":80", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getPage(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got %s request\n", r.URL.Path)
	//If it's the root, load the latest entry
	if r.URL.Path == "/" {
		loadContent(w, r, len(index.Content) - 1)
	} else {
		//Convert the URL to a string (minus the leading slash)
		s := strings.ReplaceAll(r.URL.Path, "/", "")
		i, err := strconv.Atoi(s)
		//Error Checking?
		if err != nil {}
		if i > len(index.Content) {
			return
		}
		loadContent(w, r, i)
	}
}

func loadContent(w http.ResponseWriter, r *http.Request, i int) {
	var data []byte
	var err error
	if (index.Content[i].ExecutionMethod == "execute"){
		//Run the file, and print it's output
		dir, err := filepath.Abs("Data/" + index.Content[i].Location)
		if err != nil {
            log.Fatal(err)
    	}

		data, err = exec.Command(dir, getExecutionData(index.Content[i].ExecutionData, r)...).Output()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		//Load the file from a location
		data, err = os.ReadFile("Data/" + index.Content[i].Location)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Headers which can be loaded from JSON
	if index.Content[i].ContentType != "" { w.Header().Set("Content-Type", index.Content[i].ContentType) }
	if index.Content[i].ContentDisposition != "" { w.Header().Set("Content-Disposition", index.Content[i].ContentDisposition) }
	w.Write(data)
}

func getExecutionData(requestedData []string, r *http.Request) []string {
	var data []string 
	for i := 0; i < len(requestedData); i++ {
		switch requestedData[i] {
		case "addr":
			data = append(data, r.RemoteAddr)
		}
	}
	return data
}