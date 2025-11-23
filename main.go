package main

import (
	"fmt"
	"os"
)

const HELP_STRING = "pass the file name as arugement to encode or decode\n" +
	"Usage: huffman zip|unzip -i input_file [-o output_file]\n" +
	"  zip        : encode\n" +
	"  unzip      : decode\n" +
	"  -i         : specify input file name\n" +
	"  -o         : specify output file name (optional)\n" +
	"  help, -h   : display this help message"

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

	var inputFileName string
	var outputFileName string

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
			inputFileName = os.Args[index+1]
			index++

		case "-o":
			if index == len(os.Args)-1 {
				fmt.Println("Error: -o need argument")
				os.Exit(1)
			}
			outputFileName = os.Args[index+1]
			index++

		default:
			fmt.Printf("Error: unknown argument %s\n", os.Args[index])
			os.Exit(1)
		}
		index++
	}

	if outputFileName == "" {
		if encode_flag {
			outputFileName = "out.bin"
		} else if decode_flag {
			outputFileName = "out.txt"
		}
	}

	if inputFileName == "" {
		fmt.Println("Error: input file required")
		os.Exit(1)
	}

	str, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Error: can't open file %v", inputFileName)
		os.Exit(1)
	}

	if encode_flag {
		var huffmancodes HuffmanCodes
		huffmancodes, err := GetHuffmanCodes(string(str))
		if err != nil {
			fmt.Printf("Error: generate huffman table faild:\n%v\n", err.Error())
			os.Exit(1)
		}

		// open output file
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Printf("Error: can't open output file %s\n", outputFileName)
			os.Exit(1)
		}
		defer outputFile.Close()

		// write to output
		var huffmanTableSize int
		var dataSize int
		huffmanTableSize, dataSize, err = WriteEncodeToFile(outputFile, string(str), huffmancodes)
		if err != nil {
			fmt.Printf("Error: failed to write encoded data:\n%v\n", err)
			os.Exit(1)
		}

		// print statistic information
		fmt.Printf("Original length: %d bytes\n", len(str))
		fmt.Printf("Huffman table size: %d bytes\n", huffmanTableSize)
		fmt.Printf("Compressed length (data only): %d bytes)\n", dataSize)
		fmt.Printf("Compressed length (with Huffman table): %d bytes\n", huffmanTableSize+dataSize)
		if len(str) > 0 {
			ratio := float64(huffmanTableSize+dataSize) / float64(len(str))
			fmt.Printf("Compression ratio: %.2f%%\n", ratio*100)
		}
	}

	if decode_flag {
		txt := HuffmanDecode(str)
		// handle corrupted file
		// empty text after encode is 11 bytes header
		if txt == "" && len(str) > 11 {
			fmt.Printf("Error: file %s is corrupted or not a valid huffman encoded file\n", inputFileName)
			os.Exit(1)
		}
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Printf("Error: can't open output file %s\n", outputFileName)
			os.Exit(1)
		}
		defer outputFile.Close()
		_, err = outputFile.WriteString(txt)
		if err != nil {
			fmt.Printf("Error: failed to write decoded data: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Decoded successfully, result in: %s\n", outputFileName)
	}
}
