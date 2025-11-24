package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// try to open output file, create directory if not exist
func OpenFile(filePath string) (file *os.File, err error) {
	// if output directory not exist, create it
	var dirPath string = filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("create output directory %s failed: %v", dirPath, err.Error())
		}
	}

	// create output file
	file, err = os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create output file %s failed: %v", filePath, err.Error())
	}
	return file, nil
}

// get all files in directory
func GetFilesInDir(dirPath string) (filePaths []string, batchErrors []BatchError, err error) {
	filePaths = make([]string, 0)
	// walk through directory
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			batchErrors = append(batchErrors, BatchError{Path: path, Err: walkErr})
			return nil // continue
		}
		if info.IsDir() {
			return nil
		}
		filePaths = append(filePaths, path)
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("walk through directory %s failed: %v", dirPath, err.Error())
	}
	return filePaths, batchErrors, nil
}

// get output paths for input files
// replace extension with given extension
func GetOutputPaths(inputPaths []string, outputDir string, extension string) (outputPaths []string, errors []BatchError) {
	outputPaths = make([]string, 0)
	for _, inputPath := range inputPaths {
		// replace extension with given extension
		var base string = filepath.Base(inputPath)
		var ext string = filepath.Ext(base)
		if ext != "" {
			base = base[:len(base)-len(ext)]
		}
		var outputPath string = filepath.Join(outputDir, base+"."+extension)
		outputPaths = append(outputPaths, outputPath)
	}
	return outputPaths, errors
}
