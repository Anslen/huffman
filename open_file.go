package main

import (
	"fmt"
	"os"
	"path"
)

// try to open output file, create directory if not exist
func OpenFile(filePath string) (file *os.File, err error) {
	// if output directory not exist, create it
	var dirPath string = path.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("create output directory %s failed:\n%v", dirPath, err.Error())
		}
	}

	// create output file
	file, err = os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create output file %s failed:\n%v", filePath, err.Error())
	}
	return file, nil
}
