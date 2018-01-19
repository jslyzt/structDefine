package mark

// BitMark 比特标记
type BitMark struct {
	data []uint8
}

// Set 设置标记
func (mark *BitMark) Set(id uint32, bset bool) {
	if id <= 0 {
		return
	}
	id--
	findex := int(id) / 8
	sindex := uint8(id % 8)
	ldata := len(mark.data)
	if findex+1 > ldata {
		if bset == false {
			return
		}
		ndata := make([]uint8, findex+1)
		for index := 0; index < ldata; index++ {
			ndata[index] = mark.data[index]
		}
		mark.data = ndata
	}
	mark.set(&mark.data[findex], sindex, bset)
}

func (mark *BitMark) set(data *uint8, index uint8, bset bool) {
	if data == nil || index >= 8 {
		return
	}
	if bset == true {
		*data |= uint8(1 << index)
	} else {
		*data ^= uint8(1 << index)
	}
}

// Get 获取标记
func (mark *BitMark) Get(id uint32) bool {
	if id <= 0 || mark.data == nil {
		return false
	}
	id--
	findex := int(id) / 8
	if findex >= len(mark.data) {
		return false
	}
	sindex := uint8(id % 8)
	return (mark.data[findex] & uint8(1<<sindex)) != 0
}

// Is 是否标记
func (mark *BitMark) Is(id uint32) bool {
	return mark.Get(id)
}

// Clear 清除所有标记
func (mark *BitMark) Clear() {
	mark.data = nil
}
