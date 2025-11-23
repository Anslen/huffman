package main

// read only
//
// don't allow external modification of Result and Width
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

// write single bit to recorder
func (recorder *BitsRecorder) AddBit(value uint8) {
	if recorder.Width%8 == 0 {
		recorder.Result = append(recorder.Result, value<<7)
	} else {
		recorder.Result[recorder.Width/8] |= value << (7 - recorder.Width%8)
	}
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
	for recorder.Width%8 != 0 && valueWidth > 0 {
		var index uint8 = valueWidth - 1
		// get bit and add to recorder
		recorder.AddBit(uint8((value & (1 << index)) >> index))
	}

	// add whole byte
	for valueWidth >= 8 {
		var offset uint8 = valueWidth - 8
		var mask uint64 = 0xFF << offset
		var newValue byte = byte((value & mask) >> uint64(offset))
		recorder.Result = append(recorder.Result, newValue)
	}

	// add left bits
	for valueWidth > 0 {
		var index uint8 = valueWidth - 1
		// get bit and add to recorder
		recorder.AddBit(uint8((value & (1 << index)) >> index))
	}
}
