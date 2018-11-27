package fatfs

type RawDirEntry [32]byte

func (r *RawDirEntry) Name() []byte {
	return r[0 : 0+11]
}

func (r *RawDirEntry) RawAttr() []byte {
	return r[11 : 11+1]
}

func (r *RawDirEntry) Attr() uint8 {
	return r[11]
}

func (r *RawDirEntry) SetAttr(x uint8) {
	r[11] = x
}

func (r *RawDirEntry) RawCrtTimeTenth() []byte {
	return r[13 : 13+1]
}

func (r *RawDirEntry) CrtTimeTenth() uint8 {
	return r[13]
}

func (r *RawDirEntry) SetCrtTimeTenth(x uint8) {
	r[13] = x
}

func (r *RawDirEntry) RawCrtTime() []byte {
	return r[14 : 14+2]
}

func (r *RawDirEntry) CrtTime() uint16 {
	return uint16(r[14]) | (uint16(r[14+1]) << 8)
}

func (r *RawDirEntry) SetCrtTime(x uint16) {
	r[14] = uint8(x)
	r[14+1] = uint8(x >> 8)
}

func (r *RawDirEntry) RawCrtDate() []byte {
	return r[16 : 16+2]
}

func (r *RawDirEntry) CrtDate() uint16 {
	return uint16(r[16]) | (uint16(r[16+1]) << 8)
}

func (r *RawDirEntry) SetCrtDate(x uint16) {
	r[16] = uint8(x)
	r[16+1] = uint8(x >> 8)
}

func (r *RawDirEntry) RawLstAccDate() []byte {
	return r[18 : 18+2]
}

func (r *RawDirEntry) LstAccDate() uint16 {
	return uint16(r[18]) | (uint16(r[18+1]) << 8)
}

func (r *RawDirEntry) SetLstAccDate(x uint16) {
	r[18] = uint8(x)
	r[18+1] = uint8(x >> 8)
}

func (r *RawDirEntry) RawFstClusHI() []byte {
	return r[20 : 20+2]
}

func (r *RawDirEntry) FstClusHI() uint16 {
	return uint16(r[20]) | (uint16(r[20+1]) << 8)
}

func (r *RawDirEntry) SetFstClusHI(x uint16) {
	r[20] = uint8(x)
	r[20+1] = uint8(x >> 8)
}

func (r *RawDirEntry) RawWrtTime() []byte {
	return r[22 : 22+2]
}

func (r *RawDirEntry) WrtTime() uint16 {
	return uint16(r[22]) | (uint16(r[22+1]) << 8)
}

func (r *RawDirEntry) SetWrtTime(x uint16) {
	r[22] = uint8(x)
	r[22+1] = uint8(x >> 8)
}

func (r *RawDirEntry) RawWrtDate() []byte {
	return r[24 : 24+2]
}

func (r *RawDirEntry) WrtDate() uint16 {
	return uint16(r[24]) | (uint16(r[24+1]) << 8)
}

func (r *RawDirEntry) SetWrtDate(x uint16) {
	r[24] = uint8(x)
	r[24+1] = uint8(x >> 8)
}

func (r *RawDirEntry) RawFstClusLO() []byte {
	return r[26 : 26+2]
}

func (r *RawDirEntry) FstClusLO() uint16 {
	return uint16(r[26]) | (uint16(r[26+1]) << 8)
}

func (r *RawDirEntry) SetFstClusLO(x uint16) {
	r[26] = uint8(x)
	r[26+1] = uint8(x >> 8)
}

func (r *RawDirEntry) RawFileSize() []byte {
	return r[28 : 28+4]
}

func (r *RawDirEntry) FileSize() uint32 {
	return uint32(r[28]) | (uint32(r[28+1]) << 8) |
		(uint32(r[28+2]) << 16) | (uint32(r[28+3]) << 24)
}

func (r *RawDirEntry) SetFileSize(x uint32) {
	r[28] = uint8(x)
	r[28+1] = uint8(x >> 8)
	r[28+2] = uint8(x >> 16)
	r[28+3] = uint8(x >> 24)
}
