package mark

// BitMark 比特标记
type BitMark struct {
	data []uint8
}

// Set 设置标记
func (mark *BitMark) Set(id uint32, set bool) {

}

// Get 获取标记
func (mark *BitMark) Get(id uint32) bool {
	return false
}

// Is 是否标记
func (mark *BitMark) Is(id uint32) bool {
	return mark.Get(id)
}

// Clear 清除所有标记
func (mark *BitMark) Clear() {
	mark.data = nil
}
