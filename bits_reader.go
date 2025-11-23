package main

type BitsReader struct {
	data           []byte
	width          int
	currentPointer int
}

// create a new bits reader
func NewBitsReader(data []byte, width int) (ret *BitsReader) {
	if data == nil {
		return ret
	}
	if width > len(data)*8 {
		return ret
	}

	ret = new(BitsReader)
	ret.data = data
	ret.width = width
	ret.currentPointer = 0
	return ret
}

// create a new bits reader from a single byte
func NewBitsReaderFromByte(data byte, width int) (ret *BitsReader) {
	if width > 8 {
		return ret
	}

	ret = new(BitsReader)
	ret.data = make([]byte, 1)
	ret.data[0] = data
	ret.width = width
	ret.currentPointer = 0
	return ret
}

// create a new bits reader from a uint64 value
func NewBitsReaderFromUint64(data uint64, width int) (ret *BitsReader) {
	if width < 0 || width > 64 {
		return ret
	}

	slice := make([]byte, 8)
	for i := 0; i < 8; i++ {
		slice[7-i] = byte(data & 0xFF)
		data >>= 8
	}
	reader := NewBitsReader(slice, 64)
	reader.Seek(64 - width)

	recorder := NewBitsRecorder()
	for i := 0; i < width; i++ {
		bit, _ := reader.GetBit()
		recorder.Add(uint64(bit), 1)
	}
	return NewBitsReader(recorder.Result, recorder.Width)
}

// seek to the given offset
func (reader *BitsReader) Seek(offset int) {
	reader.currentPointer += offset
	if reader.currentPointer < 0 {
		reader.currentPointer = 0
	} else if reader.currentPointer > reader.width {
		reader.currentPointer = reader.width - 1
	}
}

// get a single bit from the reader, return false if reach the end
func (reader *BitsReader) GetBit() (ret int8, ok bool) {
	if reader.currentPointer == reader.width {
		return 0, false
	}
	index := reader.currentPointer / 8
	offset := reader.currentPointer % 8
	mask := uint8(1) << (7 - offset)
	if reader.data[index]&mask == 0 {
		ret = 0
	} else {
		ret = 1
	}
	reader.currentPointer++
	return ret, true
}

// get an uint8 from the reader, return false if not enough bits
func (reader *BitsReader) GetUint8() (ret uint8, ok bool) {
	if reader.width-reader.currentPointer < 8 {
		return 0, false
	}
	for i := 0; i < 8; i++ {
		ret <<= 1
		bit, _ := reader.GetBit()
		ret |= uint8(bit)
	}
	return ret, true
}

// get an int8 from the reader, return false if not enough bits
func (reader *BitsReader) GetInt8() (ret int8, ok bool) {
	if reader.width-reader.currentPointer < 8 {
		return 0, false
	}
	for i := 0; i < 8; i++ {
		ret <<= 1
		bit, _ := reader.GetBit()
		ret |= int8(bit)
	}
	return ret, true
}

// get a byte from the reader, return false if not enough bits
func (reader *BitsReader) GetByte() (ret byte, ok bool) {
	ret, ok = reader.GetUint8()
	return ret, ok
}

// get n bits from the reader as uint64, return false if not enough bits
func (reader *BitsReader) GetNBits(n int) (ret uint64, ok bool) {
	if n < 0 || n > 64 {
		return 0, false
	}
	for i := 0; i < n; i++ {
		ret <<= 1
		bit, _ := reader.GetBit()
		ret |= uint64(bit)
	}
	return ret, true
}

// get an uint64 from the reader, return false if not enough bits
func (reader *BitsReader) GetUint64() (ret uint64, ok bool) {
	if reader.width-reader.currentPointer < 64 {
		return 0, false
	}
	for i := 0; i < 64; i++ {
		ret <<= 1
		bit, _ := reader.GetBit()
		ret |= uint64(bit)
	}
	return ret, true
}

// get an int64 from the reader, return false if not enough bits
func (reader *BitsReader) GetInt64() (ret int64, ok bool) {
	if reader.width-reader.currentPointer < 64 {
		return 0, false
	}
	for i := 0; i < 64; i++ {
		ret <<= 1
		bit, _ := reader.GetBit()
		ret |= int64(bit)
	}
	return ret, true
}
