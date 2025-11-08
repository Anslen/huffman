package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Error: no command")
		os.Exit(1)
	}

	var encode_flag bool
	var decode_flag bool
	switch os.Args[1] {
	case "encode":
		encode_flag = true
	case "decode":
		decode_flag = true
	default:
		fmt.Printf("Error: unkonwn command %v", os.Args[1])
	}

	inputFileName := os.Args[2]
	//inputFileName := "input.txt"
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

	if encode_flag {
		str, err := os.ReadFile(inputFileName)
		if err != nil {
			fmt.Printf("Error: can't open file %v", inputFileName)
			os.Exit(1)
		}

		codes, encode, width, ok := Huffman(string(str))
		if !ok {
			fmt.Println("Error: overflow, number of words greater than 2^63")
			os.Exit(1)
		}

		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Printf("Error: can't create file %s", outputFileName)
			os.Exit(1)
		}
		writeResult(outputFile, codes, width, encode)
		outputFile.Close()
	}

	if decode_flag {

	}
}

func writeResult(file io.Writer, codes HuffmanCodes, width int64, encode []byte) {
	err := writeCodes(file, codes)
	if err != nil {
		panic(err)
	}

	err = writeEncode(file, encode, width)
	if err != nil {
		panic(err)
	}
}

func writeCodes(file io.Writer, codes HuffmanCodes) (err error) {
	// write all codes
	for key, value := range codes {
		err = binary.Write(file, binary.BigEndian, key)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.BigEndian, value.code)
		if err != nil {
			return err
		}
		err = binary.Write(file, binary.BigEndian, value.width)
		if err != nil {
			return err
		}
	}

	// write zero
	err = binary.Write(file, binary.BigEndian, byte(0))
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.BigEndian, uint8(0))
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.BigEndian, int8(0))
	if err != nil {
		return err
	}
	return nil
}

func writeEncode(file io.Writer, result []byte, width int64) (err error) {
	err = binary.Write(file, binary.BigEndian, width)
	if err != nil {
		return err
	}
	_, err = file.Write(result)
	return err
}
