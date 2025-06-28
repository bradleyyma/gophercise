package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/bradleyyma/gophercise/gopher4/parse"
)

func main() {
	// Define command-line flags
	file := flag.String("file", "gopher.json", "JSON file containing the story")
	flag.Parse()

	// Open the file specified by the user
	f, err := os.Open(*file)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}

	defer f.Close()
	// Parse the file using the Parse functio

	parse.Parse(f)

}
