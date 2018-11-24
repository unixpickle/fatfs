package fatfs

import "github.com/unixpickle/essentials"

// A File Allocation Table.
type FAT struct {
	Device BlockDevice
	Offset uint32
}

func (f *FAT) ReadFAT(dataIndex uint32) (uint32, error) {
	sector, byteIdx := f.indices(dataIndex)
	block, err := f.Device.ReadSector(sector)
	if err != nil {
		return 0, essentials.AddCtx("ReadFAT", err)
	}
	return Endian.Uint32(block[byteIdx : byteIdx+4]), nil
}

func (f *FAT) WriteFAT(dataIndex uint32, contents uint32) (err error) {
	defer essentials.AddCtxTo("WriteFAT", &err)
	sector, byteIdx := f.indices(dataIndex)
	block, err := f.Device.ReadSector(sector)
	if err != nil {
		return err
	}
	Endian.PutUint32(block[byteIdx:byteIdx+4], contents)
	return f.Device.WriteSector(sector, block)
}

func (f *FAT) indices(dataIndex uint32) (uint32, int) {
	sector := f.Offset + (dataIndex / 128)
	sectorIdx := dataIndex % 128
	return sector, int(sectorIdx) * 4
}

type FATSet []*FAT

// FindFATs locates the FAT structures on disk.
func FindFATSet(b *BootSector, d BlockDevice) FATSet {
	var res FATSet
	offset := uint32(b.RsvdSecCnt())
	for i := 0; i < int(b.NumFATs()); i++ {
		res = append(res, &FAT{
			Device: d,
			Offset: offset,
		})
	}
	return res
}

func (f FATSet) ReadFAT(dataIndex uint32) (uint32, error) {
	return f[0].ReadFAT(dataIndex)
}

func (f FATSet) WriteFAT(dataIndex uint32, contents uint32) error {
	for _, x := range f {
		if err := x.WriteFAT(dataIndex, contents); err != nil {
			return err
		}
	}
	return nil
}
