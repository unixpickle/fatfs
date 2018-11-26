package fatfs

import (
	"errors"
	"math/rand"
	"strings"
)

// NewBootSector32 creates a BootSector for a new FAT32
// file-system.
func NewBootSector32(numSectors uint32, volumeLabel string) (*BootSector, error) {
	for len(volumeLabel) < 11 {
		volumeLabel += " "
	}
	if numSectors < 8*65525 {
		return nil, errors.New("volume is too small")
	} else if numSectors >= (1<<32 - 1) {
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
	res.SetTotSec32(numSectors)
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

func ceilDiv(num, denom uint32) uint32 {
	if num%denom != 0 {
		return num/denom + 1
	} else {
		return num / denom
	}
}
