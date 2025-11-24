package main

type BitsReader struct {
	data           []byte
	width          int
	currentPointer int
}

// create a new bits reader
func NewBitsReader(data []byte, width int) (ret *BitsReader) {
	if data == nil {
		return nil
	}
	if width > len(data)*8 {
		return nil
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
		return nil
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
		return nil
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
	return NewBitsReader(recorder.Result(), recorder.Width())
}

// seek to the given offset
func (reader *BitsReader) Seek(offset int) {
	reader.currentPointer += offset
	if reader.currentPointer < 0 {
		reader.currentPointer = 0
	} else if reader.currentPointer > reader.width {
		reader.currentPointer = reader.width
	}
}

// get a single bit from the reader, return false if reach the end
func (reader *BitsReader) GetBit() (ret uint8, ok bool) {
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

	// if align with byte, directly return
	if reader.currentPointer%8 == 0 {
		ret = reader.data[reader.currentPointer/8]
		reader.currentPointer += 8
		return ret, true
	}

	// if not align, read each bit
	for i := 0; i < 8; i++ {
		ret <<= 1
		bit, _ := reader.GetBit()
		ret |= uint8(bit)
	}
	return ret, true
}

// get an int8 from the reader, return false if not enough bits
func (reader *BitsReader) GetInt8() (ret int8, ok bool) {
	result, ok := reader.GetUint8()
	ret = int8(result)
	return ret, ok
}

// get a byte from the reader, return false if not enough bits
func (reader *BitsReader) GetByte() (ret byte, ok bool) {
	return reader.GetUint8()
}

// get n bits from the reader as uint64, return false if not enough bits
func (reader *BitsReader) GetNBits(n int) (ret uint64, ok bool) {
	if n < 0 || n > 64 {
		return 0, false
	}
	if reader.width-reader.currentPointer < n {
		return 0, false
	}
	// if align with byte and n is multiple of 8, read whole bytes
	if reader.currentPointer%8 == 0 && n%8 == 0 {
		numBytes := n / 8
		for i := 0; i < numBytes; i++ {
			val, _ := reader.GetByte()
			ret <<= 8
			ret |= uint64(val)
		}
		return ret, true
	}

	// read each bit
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

	var leftBits = 64
	// when not align with byte, read each bit
	for reader.currentPointer%8 != 0 {
		var bit uint8
		bit, _ = reader.GetBit()
		ret <<= 1
		ret |= uint64(bit)
		leftBits--
	}

	// read whole byte
	for leftBits >= 8 {
		var val byte
		val, _ = reader.GetByte()
		ret <<= 8
		ret |= uint64(val)
		leftBits -= 8
	}

	// read left bits
	for leftBits > 0 {
		var bit uint8
		bit, _ = reader.GetBit()
		ret <<= 1
		ret |= uint64(bit)
		leftBits--
	}
	return ret, true
}

// get an int64 from the reader, return false if not enough bits
func (reader *BitsReader) GetInt64() (ret int64, ok bool) {
	result, ok := reader.GetUint64()
	ret = int64(result)
	return ret, ok
}
