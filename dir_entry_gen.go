package fatfs

type DirEntry [32]byte

func (d *DirEntry) Name() []byte {
	return d[0 : 0+11]
}

func (d *DirEntry) RawAttr() []byte {
	return d[11 : 11+1]
}

func (d *DirEntry) Attr() uint8 {
	return d[11]
}

func (d *DirEntry) SetAttr(x uint8) {
	d[11] = x
}

func (d *DirEntry) RawCrtTimeTenth() []byte {
	return d[13 : 13+1]
}

func (d *DirEntry) CrtTimeTenth() uint8 {
	return d[13]
}

func (d *DirEntry) SetCrtTimeTenth(x uint8) {
	d[13] = x
}

func (d *DirEntry) RawCrtTime() []byte {
	return d[14 : 14+2]
}

func (d *DirEntry) CrtTime() uint16 {
	return uint16(d[14]) | (uint16(d[14+1]) << 8)
}

func (d *DirEntry) SetCrtTime(x uint16) {
	d[14] = uint8(x)
	d[14+1] = uint8(x >> 8)
}

func (d *DirEntry) RawCrtDate() []byte {
	return d[16 : 16+2]
}

func (d *DirEntry) CrtDate() uint16 {
	return uint16(d[16]) | (uint16(d[16+1]) << 8)
}

func (d *DirEntry) SetCrtDate(x uint16) {
	d[16] = uint8(x)
	d[16+1] = uint8(x >> 8)
}

func (d *DirEntry) RawLstAccDate() []byte {
	return d[18 : 18+2]
}

func (d *DirEntry) LstAccDate() uint16 {
	return uint16(d[18]) | (uint16(d[18+1]) << 8)
}

func (d *DirEntry) SetLstAccDate(x uint16) {
	d[18] = uint8(x)
	d[18+1] = uint8(x >> 8)
}

func (d *DirEntry) RawFstClusHI() []byte {
	return d[20 : 20+2]
}

func (d *DirEntry) FstClusHI() uint16 {
	return uint16(d[20]) | (uint16(d[20+1]) << 8)
}

func (d *DirEntry) SetFstClusHI(x uint16) {
	d[20] = uint8(x)
	d[20+1] = uint8(x >> 8)
}

func (d *DirEntry) RawWrtTime() []byte {
	return d[22 : 22+2]
}

func (d *DirEntry) WrtTime() uint16 {
	return uint16(d[22]) | (uint16(d[22+1]) << 8)
}

func (d *DirEntry) SetWrtTime(x uint16) {
	d[22] = uint8(x)
	d[22+1] = uint8(x >> 8)
}

func (d *DirEntry) RawWrtDate() []byte {
	return d[24 : 24+2]
}

func (d *DirEntry) WrtDate() uint16 {
	return uint16(d[24]) | (uint16(d[24+1]) << 8)
}

func (d *DirEntry) SetWrtDate(x uint16) {
	d[24] = uint8(x)
	d[24+1] = uint8(x >> 8)
}

func (d *DirEntry) RawFstClusLO() []byte {
	return d[26 : 26+2]
}

func (d *DirEntry) FstClusLO() uint16 {
	return uint16(d[26]) | (uint16(d[26+1]) << 8)
}

func (d *DirEntry) SetFstClusLO(x uint16) {
	d[26] = uint8(x)
	d[26+1] = uint8(x >> 8)
}

func (d *DirEntry) RawFileSize() []byte {
	return d[28 : 28+4]
}

func (d *DirEntry) FileSize() uint32 {
	return uint32(d[28]) | (uint32(d[28+1]) << 8) |
		(uint32(d[28+2]) << 16) | (uint32(d[28+3]) << 24)
}

func (d *DirEntry) SetFileSize(x uint32) {
	d[28] = uint8(x)
	d[28+1] = uint8(x >> 8)
	d[28+2] = uint8(x >> 16)
	d[28+3] = uint8(x >> 24)
}
