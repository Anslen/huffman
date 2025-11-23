package main

import (
	"fmt"
	"io"
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
		codes, codeWidth, result, resultWidth, ok := HuffmanEncode(string(str))
		if !ok {
			fmt.Println("Error: huffman code longer than 64 bits")
			os.Exit(1)
		}

		// calculate compression ratio
		originalBytes := len(str)
		compressedDataBytes := int((resultWidth + 7) / 8)
		codeTableBytes := (len(codes)+1)*int(2+codeWidth/8) + 1 // each code: 1 byte char + codeWidth/8 bytes code + 1 byte width
		// 8 empty bytes when storing width
		compressedWithCodesBytes := codeTableBytes + 8 + compressedDataBytes

		fmt.Printf("Original length: %d bytes\n", originalBytes)
		fmt.Printf("Compressed length (data only): %d bytes (%d bits)\n", compressedDataBytes, resultWidth)
		fmt.Printf("Compressed length (with Huffman table): %d bytes\n", compressedWithCodesBytes)
		if originalBytes > 0 {
			ratio := float64(compressedWithCodesBytes) / float64(originalBytes)
			fmt.Printf("Compression ratio: %.2f%%\n", ratio*100)
		}

		// store result
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Printf("Error: can't open output file %s\n", outputFileName)
			os.Exit(1)
		}
		defer outputFile.Close()
		if err := writeResult(outputFile, codes, codeWidth, resultWidth, result); err != nil {
			fmt.Printf("Error: failed to write encoded data: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Encoded successfully, result in: %s\n", outputFileName)
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

func writeResult(file io.Writer, codes HuffmanCodes, codeWidth uint8, width int64, encode []byte) error {
	err := writeCodes(file, codes, codeWidth)
	if err != nil {
		return fmt.Errorf("failed to write huffman codes: %w", err)
	}

	// write zero and width info
	recorder := NewBitsRecorder()
	recorder.Add(uint64(0), 16+codeWidth)
	recorder.Add(uint64(width), 64)
	_, err = file.Write(recorder.Result)
	if err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	_, err = file.Write(encode)
	if err != nil {
		return fmt.Errorf("failed to write encoded data: %w", err)
	}
	return nil
}

func writeCodes(file io.Writer, codes HuffmanCodes, codeWidth uint8) error {
	// write all codes
	recorder := NewBitsRecorder()
	recorder.Add(uint64(codeWidth), 8)
	for key, value := range codes {
		recorder.Add(uint64(key), 8)
		recorder.Add(value.Code, codeWidth)
		recorder.Add(uint64(value.Width), 8)
	}

	_, err := file.Write(recorder.Result)
	if err != nil {
		return fmt.Errorf("failed to write codes table: %w", err)
	}
	return nil
}
