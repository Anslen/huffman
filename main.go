package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func main() {
	/*
		if len(os.Args) == 1 {
			fmt.Println("Error: no command")
			os.Exit(1)
		}
	*/

	//inputFileName := os.Args[1]
	inputFileName := "input.txt"
	outputFileName := "out.bin"

	for index, command := range os.Args {
		if command == "-o" {
			if index == len(os.Args)-1 {
				fmt.Println("Error: -o need argument")
				os.Exit(1)
			}
			outputFileName = os.Args[index+1]
		}
	}

	str, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Printf("Error: can't open file %v", inputFileName)
		os.Exit(1)
	}

	codes, result, width, ok := Huffman(string(str))
	if !ok {
		fmt.Println("Error: overflow, number of words greater than 2^63")
		os.Exit(1)
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Printf("Error: can't create file %s", outputFileName)
		os.Exit(1)
	}

	err = writeCodes(outputFile, codes)
	if err != nil {
		panic(err)
	}

	err = writeResult(outputFile, result, width)
	if err != nil {
		panic(err)
	}
}

func writeCodes(file io.Writer, codes HuffmanCodes) (err error) {
	for key, value := range codes {
		err = binary.Write(file, binary.BigEndian, key)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.BigEndian, value.First)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.BigEndian, value.Second)
		if err != nil {
			return err
		}
	}
	err = binary.Write(file, binary.BigEndian, rune(0))
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.BigEndian, uint64(0))
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.BigEndian, int8(0))
	if err != nil {
		return err
	}
	return nil
}

func writeResult(file io.Writer, result []byte, width int64) (err error) {
	err = binary.Write(file, binary.BigEndian, width)
	if err != nil {
		return err
	}
	_, err = file.Write(result)
	return err
}
