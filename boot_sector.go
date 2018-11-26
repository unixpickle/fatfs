package fatfs

import (
	"errors"
	"math/rand"
	"strings"
)

// BootSector contains the data in the first sector of the
// volume.
type BootSector Sector

// NewBootSector32 creates a BootSector for a new FAT32
// file-system.
func NewBootSector32(volumeSize uint64, volumeLabel string) (*BootSector, error) {
	for len(volumeLabel) < 11 {
		volumeLabel += " "
	}
	if volumeSize < SectorSize*65525 {
		return nil, errors.New("volume is too small")
	} else if volumeSize >= SectorSize*(1<<32) {
		return nil, errors.New("volume is too large")
	}
	res := new(BootSector)
	copy(res.BootJump(), []byte{0xeb, 0, 0x90})
	copy(res.OEMName(), []byte("MSWIN4.1"))
	res.SetBytesPerSec(SectorSize)
	res.SetSecPerClus(8)
	res.SetRsvdSecCnt(2)
	res.SetNumFATs(2)
	res.SetRootEntCnt(0)
	res.SetTotSec16(0)
	res.SetMedia(0xf8)
	res.SetFatSz16(0)
	res.SetSecPerTrk(1)
	res.SetNumHeads(1)
	res.SetHiddSec(0)
	res.SetTotSec32(uint32(volumeSize / SectorSize))
	res.SetFatSz32(ceilDiv(res.TotSec32(), 4096/4))
	res.SetExtFlags(0)
	res.SetFSVer(0)
	res.SetRootClus(2)
	res.SetFSInfo(1)
	res.SetBkBootSec(0)
	res.SetDrvNum(0x80)
	res.SetBootSig(0x29)
	res.SetVolID(uint32(rand.Int31()))
	copy(res.VolLab(), []byte(strings.ToUpper(volumeLabel)))
	copy(res.FilSysType(), []byte("FAT32   "))
	res[510] = 0x55
	res[511] = 0xaa
	return res, nil
}

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

func (b *BootSector) RawFatSz32() []byte {
	return b[36 : 36+4]
}

func (b *BootSector) FatSz32() uint32 {
	return uint32(b[36]) | (uint32(b[36+1]) << 8) |
		(uint32(b[36+2]) << 16) | (uint32(b[36+3]) << 24)
}

func (b *BootSector) SetFatSz32(x uint32) {
	b[36] = uint8(x)
	b[36+1] = uint8(x >> 8)
	b[36+2] = uint8(x >> 16)
	b[36+3] = uint8(x >> 24)
}

func (b *BootSector) RawExtFlags() []byte {
	return b[40 : 40+2]
}

func (b *BootSector) ExtFlags() uint16 {
	return uint16(b[40]) | (uint16(b[40+1]) << 8)
}

func (b *BootSector) SetExtFlags(x uint16) {
	b[40] = uint8(x)
	b[40+1] = uint8(x >> 8)
}

func (b *BootSector) RawFSVer() []byte {
	return b[42 : 42+2]
}

func (b *BootSector) FSVer() uint16 {
	return uint16(b[42]) | (uint16(b[42+1]) << 8)
}

func (b *BootSector) SetFSVer(x uint16) {
	b[42] = uint8(x)
	b[42+1] = uint8(x >> 8)
}

func (b *BootSector) RawRootClus() []byte {
	return b[44 : 44+4]
}

func (b *BootSector) RootClus() uint32 {
	return uint32(b[44]) | (uint32(b[44+1]) << 8) |
		(uint32(b[44+2]) << 16) | (uint32(b[44+3]) << 24)
}

func (b *BootSector) SetRootClus(x uint32) {
	b[44] = uint8(x)
	b[44+1] = uint8(x >> 8)
	b[44+2] = uint8(x >> 16)
	b[44+3] = uint8(x >> 24)
}

func (b *BootSector) RawFSInfo() []byte {
	return b[48 : 48+2]
}

func (b *BootSector) FSInfo() uint16 {
	return uint16(b[48]) | (uint16(b[48+1]) << 8)
}

func (b *BootSector) SetFSInfo(x uint16) {
	b[48] = uint8(x)
	b[48+1] = uint8(x >> 8)
}

func (b *BootSector) RawBkBootSec() []byte {
	return b[50 : 50+2]
}

func (b *BootSector) BkBootSec() uint16 {
	return uint16(b[50]) | (uint16(b[50+1]) << 8)
}

func (b *BootSector) SetBkBootSec(x uint16) {
	b[50] = uint8(x)
	b[50+1] = uint8(x >> 8)
}

func (b *BootSector) Reserved() []byte {
	return b[52 : 52+12]
}

func (b *BootSector) RawDrvNum() []byte {
	return b[64 : 64+1]
}

func (b *BootSector) DrvNum() uint8 {
	return b[64]
}

func (b *BootSector) SetDrvNum(x uint8) {
	b[64] = x
}

func (b *BootSector) RawReserved1() []byte {
	return b[65 : 65+1]
}

func (b *BootSector) Reserved1() uint8 {
	return b[65]
}

func (b *BootSector) SetReserved1(x uint8) {
	b[65] = x
}

func (b *BootSector) RawBootSig() []byte {
	return b[66 : 66+1]
}

func (b *BootSector) BootSig() uint8 {
	return b[66]
}

func (b *BootSector) SetBootSig(x uint8) {
	b[66] = x
}

func (b *BootSector) RawVolID() []byte {
	return b[67 : 67+4]
}

func (b *BootSector) VolID() uint32 {
	return uint32(b[67]) | (uint32(b[67+1]) << 8) |
		(uint32(b[67+2]) << 16) | (uint32(b[67+3]) << 24)
}

func (b *BootSector) SetVolID(x uint32) {
	b[67] = uint8(x)
	b[67+1] = uint8(x >> 8)
	b[67+2] = uint8(x >> 16)
	b[67+3] = uint8(x >> 24)
}

func (b *BootSector) VolLab() []byte {
	return b[71 : 71+11]
}

func (b *BootSector) FilSysType() []byte {
	return b[82 : 82+8]
}

func ceilDiv(num, denom uint32) uint32 {
	if num%denom != 0 {
		return num/denom + 1
	} else {
		return num / denom
	}
}
