package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	fileListPath := os.Args[1]

	fmt.Printf("file list %s: ", fileListPath)

	var arr []string
	err := json.Unmarshal([]byte(fileListPath), &arr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("json %v: ", arr)

	// content, err := os.ReadFile("./tests/sample_invalid_satd_1.go")
	content, err := os.ReadFile(arr[0])
	if err != nil {
		panic(err)
	}

	satds, err := Parse(string(content))
	if err != nil {
		panic(err)
	}

	if len(satds) == 0 {
		fmt.Println("No valid SATDs in file")
	}

	for _, satd := range satds {
		fmt.Println(satd.Type)
		fmt.Println(satd.Description)
		fmt.Println(satd.Line)
	}
}
