package fatfs

const SectorSize = 512

type Sector [SectorSize]byte

// A BlockDevice is a raw interface for reading and
// writing blocks of data.
type BlockDevice interface {
	ReadSector(idx uint32) (*Sector, error)
	WriteSector(idx uint32, value *Sector) error
}
