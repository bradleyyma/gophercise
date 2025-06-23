package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bradleyyma/gophercise/gopher3/cyoa"
)

func main() {
	file := flag.String("file", "gopher.json", "JSON file containing the story")
	port := flag.Int("port", 8080, "Port to run the web server on")
	flag.Parse()
	fmt.Printf("Using story file: %s\n", *file)

	f, err := os.Open(*file)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	story, err := cyoa.Jsonstory(f)
	if err != nil {
		fmt.Printf("Error reading story: %v\n", err)
		return
	}

	h := cyoa.NewHandler(story, cyoa.WithPathFunc(pathFn))
	mux := http.NewServeMux()
	mux.Handle("/story/", h)
	fmt.Printf("Starting server on port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := r.URL.Path
	if path == "/story" || path == "story" || path == "" || path == "/" {
		path = "/story/intro"
	}
	return path[len("/story/"):] // Remove leading slash
}
