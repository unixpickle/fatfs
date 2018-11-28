package fatfs

import (
	"errors"
	"io"
	"os"

	"github.com/unixpickle/essentials"
)

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

// FileDevice is a BlockDevice that is backed by a file,
// possibly a block device.
type FileDevice struct {
	file *os.File
	size uint32
}

// NewFileDevice creates a FileDevice that wraps a file f.
func NewFileDevice(f *os.File) (*FileDevice, error) {
	size, err := getDeviceSize(f)
	if err != nil {
		return nil, essentials.AddCtx("NewFileDevice", err)
	}
	if size >= 512*(1<<32) {
		size = 512*(1<<32) - 1
	}
	return &FileDevice{
		file: f,
		size: uint32(size / 512),
	}, nil
}

func (f *FileDevice) NumSectors() uint32 {
	return f.size
}

func (f *FileDevice) ReadSector(idx uint32) (sec *Sector, err error) {
	defer essentials.AddCtxTo("ReadSector", &err)
	if idx >= f.size {
		return nil, errors.New("sector out of bounds")
	}
	if _, err := f.file.Seek(int64(idx)*SectorSize, io.SeekStart); err != nil {
		return nil, err
	}
	var res Sector
	if _, err := io.ReadFull(f.file, res[:]); err != nil {
		return nil, err
	}
	return &res, nil
}

func (f *FileDevice) WriteSector(idx uint32, data *Sector) (err error) {
	defer essentials.AddCtxTo("WriteSector", &err)
	if idx >= f.size {
		return errors.New("sector out of bounds")
	}
	_, err = f.file.WriteAt(data[:], int64(idx)*SectorSize)
	return err
}
