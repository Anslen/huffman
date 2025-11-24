package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const HELP_STRING = "pass the file name as arugement to encode or decode\n" +
	"Usage: huffman zip|unzip -i input_file [-o output_file]\n" +
	"  zip        : encode\n" +
	"  unzip      : decode\n" +
	"  -i         : specify input file name\n" +
	"  -o         : specify output file name (optional)\n" +
	"  help, -h   : display this help message"

// processPath converts relative output path to absolute path in the same directory as input file
//
// if inputPath is absolute and outputPath is relative, place output file in input file's directory
func processPath(inputPath, outputPath string) (string, string) {

	// Clean the paths
	inputPath = filepath.Clean(inputPath)

	if filepath.IsAbs(inputPath) && !filepath.IsAbs(outputPath) {
		inputDir := filepath.Dir(inputPath)
		outputAbsPath := filepath.Clean(filepath.Join(inputDir, outputPath))
		return inputPath, outputAbsPath
	} else {
		outputPath = filepath.Clean(outputPath)
		return inputPath, outputPath
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println(HELP_STRING)
		os.Exit(0)
	}

	var encode_flag bool = os.Args[1] == "zip"
	var decode_flag bool = os.Args[1] == "unzip"

	if (!encode_flag) && (!decode_flag) {
		fmt.Println("Error: first argument must be 'zip' or 'unzip'")
		os.Exit(1)
	}

	var inputPath string
	var outputPath string

	// read arguments
	index := 2
	for index < len(os.Args) {
		switch os.Args[index] {
		case "-h", "help":
			fmt.Println(HELP_STRING)
			os.Exit(0)

		case "-i":
			if index == len(os.Args)-1 {
				fmt.Println("Error: -i need argument")
				os.Exit(1)
			}
			inputPath = os.Args[index+1]
			index++

		case "-o":
			if index == len(os.Args)-1 {
				fmt.Println("Error: -o need argument")
				os.Exit(1)
			}
			outputPath = os.Args[index+1]
			index++

		default:
			fmt.Printf("Error: unknown argument %s\n", os.Args[index])
			os.Exit(1)
		}
		index++
	}

	if outputPath == "" {
		if encode_flag {
			outputPath = "out.bin"
		} else if decode_flag {
			outputPath = "out.txt"
		}
	}

	// Convert relative output path to absolute if input is absolute
	inputPath, outputPath = processPath(inputPath, outputPath)

	if inputPath == "" {
		fmt.Println("Error: input file required")
		os.Exit(1)
	}

	if encode_flag {
		fmt.Printf("Compressing...")

		// encode file
		var encodeSize EncodeSize
		var encodeTime EncodeTime
		encodeSize, encodeTime, err := Encode(inputPath, outputPath)
		if err != nil {
			fmt.Printf("Error: write encoded data failed:\n%v\n", err)
			os.Exit(1)
		}

		// read size information
		var originalSize int = encodeSize.orininal
		var huffmanTableSize int = encodeSize.HuffmanTable
		var encodedDataSize int = encodeSize.EncodedData
		var encodedSize int = huffmanTableSize + encodedDataSize

		// read time information
		var codeGenTime time.Duration = encodeTime.CodeGenTime
		var writeTime time.Duration = encodeTime.WriteFileTime

		// print statistic information
		fmt.Printf("\nEncode successful, result in: %v\n\n", outputPath)
		fmt.Printf("Original size: %d bytes\n", originalSize)
		fmt.Printf("Huffman table size: %d bytes\n", huffmanTableSize)
		fmt.Printf("Compressed size (data only): %d bytes\n", encodedDataSize)
		fmt.Printf("Compressed size (with Huffman table): %d bytes\n", encodedSize)
		if originalSize > 0 {
			ratio := float64(encodedSize) / float64(originalSize)
			fmt.Printf("Compression ratio: %.2f%%\n\n", ratio*100)
		}
		totalTime := codeGenTime + writeTime
		fmt.Printf("Time: Huffman table generation: %.2fs, File writing: %.2fs, Total: %.2fs\n",
			float64(codeGenTime.Milliseconds())/1000,
			float64(writeTime.Milliseconds())/1000,
			float64(totalTime.Milliseconds())/1000)
	}

	if decode_flag {
		// decode file
		fmt.Printf("Decompressing...")

		var decodeSize DecodeSize
		var decodeTime time.Duration
		decodeSize, decodeTime, err := Decode(inputPath, outputPath)
		if err != nil {
			fmt.Printf("Error: failed to decode file %s:\n%v\n", inputPath, err)
			os.Exit(1)
		}

		fmt.Printf("\nDecoded successfully, result in: %s\n", outputPath)
		fmt.Printf("Original size: %d bytes\n", decodeSize.Original)
		fmt.Printf("Decompressed size: %d bytes\n", decodeSize.Decoded)
		fmt.Printf("Time: Decoding: %.2fs\n", float64(decodeTime.Milliseconds())/1000)
	}
}
