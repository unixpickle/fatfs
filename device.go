package fatfs

const SectorSize = 512

type Sector [SectorSize]byte

// A BlockDevice is a raw interface for reading and
// writing blocks of data.
type BlockDevice interface {
	NumSectors() uint32
	ReadSector(idx uint32) (*Sector, error)
	WriteSector(idx uint32, value *Sector) error
}

// A RAMDisk is a BlockDevice that is backed by a simple
// memory buffer.
type RAMDisk []byte

func (r RAMDisk) NumSectors() uint32 {
	return uint32(len(r) / 512)
}

func (r RAMDisk) ReadSector(idx uint32) (*Sector, error) {
	var sec Sector
	copy(sec[:], r[int(idx)*SectorSize:])
	return &sec, nil
}

func (r RAMDisk) WriteSector(idx uint32, value *Sector) error {
	copy(r[int(idx)*SectorSize:], value[:])
	return nil
}
