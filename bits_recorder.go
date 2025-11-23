package main

// read only
//
// don't allow external modification of Result and Width
type BitsRecorder struct {
	result []byte
	width  int
}

// create a new bits recorder
func NewBitsRecorder() (ret *BitsRecorder) {
	ret = new(BitsRecorder)
	ret.result = make([]byte, 0)
	return ret
}

// write single bit to recorder
func (recorder *BitsRecorder) AddBit(value uint8) {
	value &= 0x1
	// extend slice if needed
	if recorder.width%8 == 0 {
		recorder.result = append(recorder.result, value<<7)
	} else {
		recorder.result[recorder.width/8] |= value << (7 - recorder.width%8)
	}
	recorder.width++
}

// add bits to the recorder
//
// low bits in value valid
//
// write in MSB-first order
func (recorder *BitsRecorder) Add(value uint64, valueWidth uint8) {
	if valueWidth == 0 || valueWidth > 64 {
		return
	}

	// add single bit when not align with 8
	for recorder.width%8 != 0 && valueWidth > 0 {
		var index uint8 = valueWidth - 1
		// get bit and add to recorder
		recorder.AddBit(uint8((value & (1 << index)) >> index))
		valueWidth--
	}

	// add whole byte
	for valueWidth >= 8 {
		var offset uint8 = valueWidth - 8
		var mask uint64 = 0xFF << offset
		var newValue byte = byte((value & mask) >> uint64(offset))
		recorder.result = append(recorder.result, newValue)
		recorder.width += 8
		valueWidth -= 8
	}

	// add left bits
	for valueWidth > 0 {
		var index uint8 = valueWidth - 1
		// get bit and add to recorder
		recorder.AddBit(uint8((value & (1 << index)) >> index))
		valueWidth--
	}
}

// Result returns the recorded bits as a byte slice
// The returned slice shares the underlying storage; do not modify it
// Invalid bits in the last byte are zeroed.
func (recorder *BitsRecorder) Result() []byte {
	return recorder.result
}

// Width returns the total number of bits recorded
func (recorder *BitsRecorder) Width() int {
	return recorder.width
}
