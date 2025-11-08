package main

type BitsRecorder struct {
	Result []byte
	Width  int
}

func NewBitsRecorder() (ret *BitsRecorder) {
	ret = new(BitsRecorder)
	ret.Result = make([]byte, 0)
	return ret
}

func (recorder *BitsRecorder) Add(value uint64, valueWidth int8) {
	var mask uint64
	var newData byte
	for valueWidth != 0 {
		if recorder.Width%8 != 0 {
			leftWidth := int8(8 - (recorder.Width % 8))
			length := len(recorder.Result)
			if leftWidth >= valueWidth {
				mask = uint64(0xff >> (8 - valueWidth))
				newData = byte((value & mask) << (leftWidth - valueWidth))
				recorder.Result[length-1] = recorder.Result[length-1] | newData
				recorder.Width += int(valueWidth)
				valueWidth = 0
			} else {
				mask = uint64(0xff >> (8 - leftWidth) << (valueWidth - leftWidth))
				newData = byte((value & mask) >> (valueWidth - leftWidth))
				recorder.Result[length-1] = recorder.Result[length-1] | newData
				recorder.Width += int(leftWidth)
				valueWidth -= leftWidth
			}
		} else {
			if valueWidth >= 8 {
				mask = uint64(0xff << (valueWidth - 8))
				newData = byte((value & mask) >> (valueWidth - 8))
				recorder.Result = append(recorder.Result, newData)
				recorder.Width += 8
				valueWidth -= 8
			} else {
				mask = uint64(0xff >> (8 - valueWidth))
				newData = byte((value & mask) << (8 - valueWidth))
				recorder.Result = append(recorder.Result, newData)
				recorder.Width += int(valueWidth)
				valueWidth = 0
			}
		}
	}
}
