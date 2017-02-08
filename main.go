package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify file to convert from")
		os.Exit(-1)
	}

	filename := os.Args[1]

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File do not exist", filename)
		os.Exit(-1)
	}

	filename, err := filepath.Abs(filename)
	if err != nil {
		fmt.Println("Could not get absolute path", err)
		os.Exit(-1)
	}

	ext := filepath.Ext(filename)
	switch ext {
	case ".yml", ".yaml":
		convertFromYaml(filename)
	case ".conf":
		convertFromConf(filename)
	default:
		fmt.Println("Unknown file extension", ext, "file is", filename)
		os.Exit(-1)
	}
}
