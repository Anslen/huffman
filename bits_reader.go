package main

type BitsReader struct {
	data           []byte
	width          int
	currentPointer int
}

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

func NewBitsReaderFromUint64(data uint64, width int) (ret *BitsReader) {
	if width > 64 {
		return ret
	}

	ret = new(BitsReader)
	// read data to byte slice
	ret.data = make([]byte, 8)
	for i := 0; i < 8; i++ {
		ret.data[7-i] = byte(data & 0xFF)
		data >>= 8
	}

	ret.width = width
	ret.currentPointer = 0
	return ret
}

func (reader *BitsReader) Seek(offset int) {
	reader.currentPointer += offset
}

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

func (reader *BitsReader) GetByte() (ret byte, ok bool) {
	ret, ok = reader.GetUint8()
	return ret, ok
}

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
