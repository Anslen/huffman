package main

type BitsRecorder struct {
	Result []byte
	Width  int
}

// create a new bits recorder
func NewBitsRecorder() (ret *BitsRecorder) {
	ret = new(BitsRecorder)
	ret.Result = make([]byte, 0)
	return ret
}

// add bits to the recorder
//
// write in MSB-first order
func (recorder *BitsRecorder) Add(value uint64, valueWidth uint8) {
	// write valueWidth bits of value from little with MSB first
	if valueWidth == 0 || valueWidth > 64 {
		return
	}

	// write bits MSB-first within the given valueWidth
	for v := int(valueWidth); v > 0; v-- {
		bitIndex := uint(v - 1)
		bit := (value >> bitIndex) & 1

		// ensure there's a byte to write into
		if recorder.Width%8 == 0 {
			recorder.Result = append(recorder.Result, 0)
		}

		byteIndex := len(recorder.Result) - 1
		bitOffset := 7 - (recorder.Width % 8) // MSB-first in each byte
		if bit == 1 {
			recorder.Result[byteIndex] |= byte(1 << uint(bitOffset))
		}
		recorder.Width++
	}
}
