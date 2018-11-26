package fatfs

import (
	"errors"

	"github.com/unixpickle/essentials"
)

// FS provides all the information needed to perform
// file-system operations.
type FS struct {
	Device     BlockDevice
	BootSector *BootSector

	fatSectors []uint32
}

// NewFS creates a file-system using the block device.
func NewFS(b BlockDevice) (*FS, error) {
	bsData, err := b.ReadSector(0)
	if err != nil {
		return nil, essentials.AddCtx("NewFS", err)
	}
	bs := BootSector(*bsData)
	fs := &FS{Device: b, BootSector: &bs}
	offset := uint32(bs.RsvdSecCnt())
	for i := 0; i < int(bs.NumFATs()); i++ {
		fs.fatSectors = append(fs.fatSectors, offset)
		offset += bs.FatSz32()
	}
	return fs, nil
}

// ClusterSize gets the number of bytes per cluster.
func (f *FS) ClusterSize() int {
	return int(f.BootSector.SecPerClus()) * SectorSize
}

// NumClusters gets the number of data clusters.
func (f *FS) NumClusters() uint32 {
	b := f.BootSector
	numSectors := b.TotSec32() - (b.FatSz32()*uint32(b.NumFATs()) + uint32(b.RsvdSecCnt()))
	return numSectors / uint32(b.SecPerClus())
}

// ReadFAT reads a FAT entry.
func (f *FS) ReadFAT(dataIndex uint32) (uint32, error) {
	sector, byteIdx := fatIndices(dataIndex)
	block, err := f.Device.ReadSector(f.fatSectors[0] + sector)
	if err != nil {
		return 0, essentials.AddCtx("ReadFAT", err)
	}
	return Endian.Uint32(block[byteIdx:byteIdx+4]) & 0x0fffffff, nil
}

// WriteFAT writes a FAT entry.
func (f *FS) WriteFAT(dataIndex uint32, contents uint32) error {
	sector, byteIdx := fatIndices(dataIndex)
	for _, sectorOffset := range f.fatSectors {
		block, err := f.Device.ReadSector(sector + sectorOffset)
		if err != nil {
			return essentials.AddCtx("WriteFAT", err)
		}
		oldContents := Endian.Uint32(block[byteIdx : byteIdx+4])
		newContents := (contents & 0x0fffffff) | (oldContents & 0xf0000000)
		Endian.PutUint32(block[byteIdx:byteIdx+4], newContents)
		err = f.Device.WriteSector(sector+sectorOffset, block)
		if err != nil {
			return essentials.AddCtx("WriteFAT", err)
		}
	}
	return nil
}
