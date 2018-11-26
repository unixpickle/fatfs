package fatfs

import (
	"errors"
	"io"

	"github.com/unixpickle/essentials"
)

const EOF = 0x0FFFFFF8

// A Chain is a readable, writeable, expandable piece of
// data on a file-system. It is stored as a sequence of
// clusters, joined together by the FAT.
//
// A Chain behaves like a tape. At any point, it is
// pointing to a cluster, and it may be moved back and
// forth, expanded, etc.
type Chain struct {
	fs      *FS
	cluster uint32
	prev    []uint32
}

// NewChain creates a Chain starting at a cluster.
func NewChain(fs *FS, start uint32) *Chain {
	return &Chain{fs: fs, cluster: start}
}

// RootDirChain gets a Chain for the root directory.
func RootDirChain(fs *FS) *Chain {
	return NewChain(fs, fs.BootSector.RootClus())
}

// ReadCluster reads the current cluster of the Chain.
func (c *Chain) ReadCluster() ([]byte, error) {
	res := make([]byte, 0, c.fs.ClusterSize())
	offset := c.clusterSector()
	for i := 0; i < int(c.fs.BootSector.SecPerClus()); i++ {
		sector, err := c.fs.Device.ReadSector(offset + uint32(i))
		if err != nil {
			return nil, essentials.AddCtx("ReadCluster", err)
		}
		res = append(res, sector[:]...)
	}
	return res, nil
}

// WriteCluster writes the current cluster of the chain.
func (c *Chain) WriteCluster(data []byte) (err error) {
	defer essentials.AddCtxTo("WriteCluster", &err)
	if len(data) != c.fs.ClusterSize() {
		return errors.New("incorrect cluster size")
	}
	offset := c.clusterSector()
	var chunk Sector
	for i := 0; i < int(c.fs.BootSector.SecPerClus()); i++ {
		copy(chunk[:], data[i*SectorSize:])
		if err := c.fs.Device.WriteSector(offset+uint32(i), &chunk); err != nil {
			return err
		}
	}
	return nil
}

// Seek moves around within the chain by a certain number
// of clusters (not bytes).
//
// It returns the new cluster offset in the chain.
//
// Seeking past the end of the chain is equivalent to
// seeking to the end of the chain.
func (c *Chain) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		if int64(len(c.prev)) > offset {
			c.cluster = c.prev[offset]
			c.prev = c.prev[:offset]
			return offset, nil
		} else if int64(len(c.prev)) == offset {
			return offset, nil
		}
		return c.Seek(offset-int64(len(c.prev)), io.SeekCurrent)
	} else if whence == io.SeekCurrent {
		if offset < 0 {
			if -offset > int64(len(c.prev)) {
				return 0, errors.New("Seek: went before the start of the chain")
			}
			newPrevLen := int64(len(c.prev)) + offset
			c.cluster = c.prev[newPrevLen]
			c.prev = c.prev[:newPrevLen]
			return int64(len(c.prev)), nil
		}
		for i := int64(0); i < offset; i++ {
			next, err := c.fs.ReadFAT(c.cluster)
			if err != nil {
				return 0, essentials.AddCtx("Seek", err)
			}
			if next >= EOF {
				return int64(len(c.prev)), nil
			}
			c.prev = append(c.prev, c.cluster)
			c.cluster = next
		}
		return int64(len(c.prev)), nil
	} else if whence == io.SeekEnd {
		if _, err := c.Seek(1<<32, io.SeekCurrent); err != nil {
			return 0, err
		}
		return c.Seek(offset, io.SeekCurrent)
	}
	return 0, errors.New("Seek: unknown whence")
}

// Extend adds a new cluster to the end of the chain and
// seeks to it.
func (c *Chain) Extend() (err error) {
	defer essentials.AddCtxTo("Extend", &err)
	if _, err := c.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	cluster, err := c.fs.Alloc()
	if err != nil {
		return err
	}
	if err := c.fs.WriteFAT(c.cluster, cluster); err != nil {
		c.fs.WriteFAT(cluster, 0)
		return err
	}
	c.prev = append(c.prev, c.cluster)
	c.cluster = cluster
	return nil
}

// Truncate removes the final cluster from the chain and
// seeks to the new end. Fails if the chain only contains
// one cluster.
func (c *Chain) Truncate() (err error) {
	defer essentials.AddCtxTo("Truncate", &err)
	if _, err := c.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	if len(c.prev) == 0 {
		return errors.New("no clusters to remove")
	}
	previous := c.prev[len(c.prev)-1]
	if err := c.fs.WriteFAT(previous, EOF); err != nil {
		return err
	}
	if err := c.fs.WriteFAT(c.cluster, 0); err != nil {
		return err
	}
	c.prev = c.prev[:len(c.prev)-1]
	c.cluster = previous
	return nil
}

// ReadFrom takes all the data from r and writes it to the
// end of the chain.
//
// Returns the number of bytes read from r before an error
// was encountered.
func (c *Chain) ReadFrom(r io.Reader) (n int64, err error) {
	defer essentials.AddCtxTo("WriteFrom", &err)
	if _, err := c.Seek(0, io.SeekEnd); err != nil {
		return 0, err
	}
	needsExtend := false
	for {
		buffer := make([]byte, c.fs.ClusterSize())
		m, readErr := io.ReadFull(r, buffer)
		n += int64(m)
		if readErr == io.EOF {
			break
		}

		if needsExtend {
			if err := c.Extend(); err != nil {
				return n, err
			}
		}
		needsExtend = true
		if err := c.WriteCluster(buffer); err != nil {
			return n, err
		}

		if readErr == io.ErrUnexpectedEOF {
			break
		} else if readErr != nil {
			return n, readErr
		}
	}
	return n, nil
}

func (c *Chain) clusterSector() uint32 {
	b := c.fs.BootSector
	firstData := uint32(b.RsvdSecCnt()) + uint32(b.NumFATs())*b.FatSz32()
	return firstData + (c.cluster-2)*uint32(b.SecPerClus())
}
