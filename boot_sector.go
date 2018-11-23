package fatfs

type BootSector [512]byte

func (b *BootSector) BootJump() []byte {
	return b[0 : 0+3]
}

func (b *BootSector) OEMName() []byte {
	return b[3 : 3+8]
}

func (b *BootSector) RawBytesPerSec() []byte {
	return b[11 : 11+2]
}

func (b *BootSector) BytesPerSec() uint16 {
	return uint16(b[11]) | (uint16(b[11+1]) << 8)
}

func (b *BootSector) SetBytesPerSec(x uint16) {
	b[11] = uint8(x)
	b[11+1] = uint8(x >> 8)
}

func (b *BootSector) RawSecPerClus() []byte {
	return b[13 : 13+1]
}

func (b *BootSector) SecPerClus() uint8 {
	return b[13]
}

func (b *BootSector) SetSecPerClus(x uint8) {
	b[13] = x
}

func (b *BootSector) RawRsvdSecCnt() []byte {
	return b[14 : 14+2]
}

func (b *BootSector) RsvdSecCnt() uint16 {
	return uint16(b[14]) | (uint16(b[14+1]) << 8)
}

func (b *BootSector) SetRsvdSecCnt(x uint16) {
	b[14] = uint8(x)
	b[14+1] = uint8(x >> 8)
}

func (b *BootSector) RawNumFATs() []byte {
	return b[16 : 16+1]
}

func (b *BootSector) NumFATs() uint8 {
	return b[16]
}

func (b *BootSector) SetNumFATs(x uint8) {
	b[16] = x
}

func (b *BootSector) RawRootEntCnt() []byte {
	return b[17 : 17+2]
}

func (b *BootSector) RootEntCnt() uint16 {
	return uint16(b[17]) | (uint16(b[17+1]) << 8)
}

func (b *BootSector) SetRootEntCnt(x uint16) {
	b[17] = uint8(x)
	b[17+1] = uint8(x >> 8)
}

func (b *BootSector) RawTotSec16() []byte {
	return b[19 : 19+2]
}

func (b *BootSector) TotSec16() uint16 {
	return uint16(b[19]) | (uint16(b[19+1]) << 8)
}

func (b *BootSector) SetTotSec16(x uint16) {
	b[19] = uint8(x)
	b[19+1] = uint8(x >> 8)
}

func (b *BootSector) RawMedia() []byte {
	return b[21 : 21+1]
}

func (b *BootSector) Media() uint8 {
	return b[21]
}

func (b *BootSector) SetMedia(x uint8) {
	b[21] = x
}

func (b *BootSector) RawFatSz16() []byte {
	return b[22 : 22+2]
}

func (b *BootSector) FatSz16() uint16 {
	return uint16(b[22]) | (uint16(b[22+1]) << 8)
}

func (b *BootSector) SetFatSz16(x uint16) {
	b[22] = uint8(x)
	b[22+1] = uint8(x >> 8)
}

func (b *BootSector) RawSecPerTrk() []byte {
	return b[24 : 24+2]
}

func (b *BootSector) SecPerTrk() uint16 {
	return uint16(b[24]) | (uint16(b[24+1]) << 8)
}

func (b *BootSector) SetSecPerTrk(x uint16) {
	b[24] = uint8(x)
	b[24+1] = uint8(x >> 8)
}

func (b *BootSector) RawNumHeads() []byte {
	return b[26 : 26+2]
}

func (b *BootSector) NumHeads() uint16 {
	return uint16(b[26]) | (uint16(b[26+1]) << 8)
}

func (b *BootSector) SetNumHeads(x uint16) {
	b[26] = uint8(x)
	b[26+1] = uint8(x >> 8)
}

func (b *BootSector) RawHiddSec() []byte {
	return b[28 : 28+4]
}

func (b *BootSector) HiddSec() uint32 {
	return uint32(b[28]) | (uint32(b[28+1]) << 8) |
		(uint32(b[28+2]) << 16) | (uint32(b[28+3]) << 24)
}

func (b *BootSector) SetHiddSec(x uint32) {
	b[28] = uint8(x)
	b[28+1] = uint8(x >> 8)
	b[28+2] = uint8(x >> 16)
	b[28+3] = uint8(x >> 24)
}

func (b *BootSector) RawTotSec32() []byte {
	return b[32 : 32+4]
}

func (b *BootSector) TotSec32() uint32 {
	return uint32(b[32]) | (uint32(b[32+1]) << 8) |
		(uint32(b[32+2]) << 16) | (uint32(b[32+3]) << 24)
}

func (b *BootSector) SetTotSec32(x uint32) {
	b[32] = uint8(x)
	b[32+1] = uint8(x >> 8)
	b[32+2] = uint8(x >> 16)
	b[32+3] = uint8(x >> 24)
}

func (b *BootSector) RawDrvNum() []byte {
	return b[36 : 36+1]
}

func (b *BootSector) DrvNum() uint8 {
	return b[36]
}

func (b *BootSector) SetDrvNum(x uint8) {
	b[36] = x
}

func (b *BootSector) RawReserved1() []byte {
	return b[37 : 37+1]
}

func (b *BootSector) Reserved1() uint8 {
	return b[37]
}

func (b *BootSector) SetReserved1(x uint8) {
	b[37] = x
}

func (b *BootSector) RawBootSig() []byte {
	return b[38 : 38+1]
}

func (b *BootSector) BootSig() uint8 {
	return b[38]
}

func (b *BootSector) SetBootSig(x uint8) {
	b[38] = x
}

func (b *BootSector) RawVolID() []byte {
	return b[39 : 39+4]
}

func (b *BootSector) VolID() uint32 {
	return uint32(b[39]) | (uint32(b[39+1]) << 8) |
		(uint32(b[39+2]) << 16) | (uint32(b[39+3]) << 24)
}

func (b *BootSector) SetVolID(x uint32) {
	b[39] = uint8(x)
	b[39+1] = uint8(x >> 8)
	b[39+2] = uint8(x >> 16)
	b[39+3] = uint8(x >> 24)
}

func (b *BootSector) VolLab() []byte {
	return b[43 : 43+11]
}

func (b *BootSector) FilSysType() []byte {
	return b[54 : 54+8]
}
