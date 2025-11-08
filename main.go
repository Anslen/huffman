package main

import (
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
		os.Exit(1)
	}

	inputFileName := os.Args[2]
	//inputFileName := "input.txt"
	var outputFileName string
	if encode_flag {
		outputFileName = "out.bin"
	} else if decode_flag {
		outputFileName = "out.txt"
	}

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

	if encode_flag {
		codes, codeWidth, result, resultWidth, ok := HuffmanEncode(string(str))
		if !ok {
			fmt.Println("Error: overflow, number of words greater than 2^63")
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
			fmt.Printf("Error: can't create file %s", outputFileName)
			os.Exit(1)
		}
		writeResult(outputFile, codes, codeWidth, resultWidth, result)
		outputFile.Close()
	}

	if decode_flag {
		txt := HuffmanDecode(str)
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Printf("Error: can't create file %s", outputFileName)
			os.Exit(1)
		}
		outputFile.WriteString(txt)
		outputFile.Close()
	}
}

func writeResult(file io.Writer, codes HuffmanCodes, codeWidth uint8, width int64, encode []byte) {
	err := writeCodes(file, codes, codeWidth)
	if err != nil {
		panic(err)
	}

	// write zero and width info
	recorder := NewBitsRecorder()
	recorder.Add(uint64(0), 16+codeWidth)
	recorder.Add(uint64(width), 64)
	_, err = file.Write(recorder.Result)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(encode)
	if err != nil {
		panic(err)
	}
}

func writeCodes(file io.Writer, codes HuffmanCodes, codeWidth uint8) (err error) {
	// write all codes
	recorder := NewBitsRecorder()
	recorder.Add(uint64(codeWidth), 8)
	for key, value := range codes {
		recorder.Add(uint64(key), 8)
		recorder.Add(value.Code, codeWidth)
		recorder.Add(uint64(value.Width), 8)
	}

	_, err = file.Write(recorder.Result)
	return err
}
