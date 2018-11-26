package fatfs

import (
	"errors"
	"io"

	"github.com/unixpickle/essentials"
)

// A Dir is an open handle to a directory.
type Dir struct {
	chain *Chain
}

// NewDir creates a Dir from a Chain.
// The Chain must be seeked to the start.
func NewDir(c *Chain) *Dir {
	return &Dir{chain: c}
}

// ReadDir reads the directory listings.
func (d *Dir) ReadDir() (entries []*DirEntry, err error) {
	defer essentials.AddCtxTo("ReadDir", &err)
	_, err = d.loopClusters(func() (bool, error) {
		cluster, err := d.chain.ReadCluster()
		if err != nil {
			return false, err
		}
		for i := 0; i < len(cluster); i += 32 {
			var entry DirEntry
			copy(entry[:], cluster[i:])
			if !entry.IsFree() {
				entries = append(entries, &entry)
			}
		}
		return false, nil
	})
	return
}

// AddEntry adds a directory entry.
//
// This will seek to an unspecified part of the directory.
// Thus, you should call Reset() before further ReadDir
// calls.
func (d *Dir) AddEntry(newEntry *DirEntry) (err error) {
	defer essentials.AddCtxTo("AddEntry", &err)

	added, err := d.loopClusters(func() (bool, error) {
		cluster, err := d.chain.ReadCluster()
		if err != nil {
			return false, err
		}
		for i := 0; i < len(cluster); i += 32 {
			var entry DirEntry
			copy(entry[:], cluster[i:])
			if entry.IsFree() {
				copy(cluster[i:], newEntry[:])
				return true, d.chain.WriteCluster(cluster)
			}
		}
		return false, nil
	})

	if added || err != nil {
		return err
	}

	// Create a new cluster of listings.
	if err := d.chain.Extend(); err != nil {
		return err
	}
	cluster := make([]byte, d.chain.fs.ClusterSize())
	copy(cluster, newEntry[:])
	return d.chain.WriteCluster(cluster)
}

// RemoveEntry deletes the entry for the given name.
func (d *Dir) RemoveEntry(name string) error {
	found, err := d.loopClusters(func() (bool, error) {
		cluster, err := d.chain.ReadCluster()
		if err != nil {
			return false, err
		}
		for i := 0; i < len(cluster); i += 32 {
			var entry DirEntry
			copy(entry[:], cluster[i:])
			if string(entry.Name()) == name {
				entry.Name()[0] = 0xe5
				copy(cluster[i:], entry[:])
				return true, d.chain.WriteCluster(cluster)
			}
		}
		return false, nil
	})
	if err == nil && !found {
		err = errors.New("file not found")
	}
	return essentials.AddCtx("RemoveEntry", err)
}

func (d *Dir) loopClusters(f func() (done bool, err error)) (bool, error) {
	if _, err := d.chain.Seek(0, io.SeekStart); err != nil {
		return false, err
	}
	for offset := int64(0); true; offset += 1 {
		if done, err := f(); err != nil {
			return false, err
		} else if done {
			return true, nil
		}
		if newOffset, err := d.chain.Seek(1, io.SeekCurrent); err != nil {
			return false, err
		} else if newOffset == offset {
			break
		}
	}
	return false, nil
}
