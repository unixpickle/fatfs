package fatfs

import (
	"io"

	"github.com/unixpickle/essentials"
)

// A Dir is an open handle to a directory.
type Dir struct {
	chain *Chain
	atEOF bool
}

// Reset seeks to the beginning of the listing.
func (d *Dir) Reset() error {
	_, err := d.chain.Seek(0, io.SeekStart)
	d.atEOF = false
	return essentials.AddCtx("Reset", err)
}

// ReadDir reads the next batch of directory entries and
// advances the directory handle.
// Returns io.EOF if the end of the listings were reached.
func (d *Dir) ReadDir() ([]*DirEntry, error) {
	cluster, err := d.chain.ReadCluster()
	if err != nil {
		return nil, essentials.AddCtx("ReadDir", err)
	}

	var results []*DirEntry
	for i := 0; i < len(cluster); i += 32 {
		var entry DirEntry
		copy(entry[:], cluster[i:])
		if !entry.IsFree() && !entry.IsLongName() {
			results = append(results, &entry)
		}
	}

	offset, err := d.chain.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, essentials.AddCtx("ReadDir", err)
	}
	newOffset, err := d.chain.Seek(1, io.SeekCurrent)
	if err != nil {
		return nil, essentials.AddCtx("ReadDir", err)
	}
	if newOffset == offset {
		d.atEOF = true
	}
	return results, nil
}

// AddEntry adds a directory entry.
//
// This will seek to an unspecified part of the directory.
// Thus, you should call Reset() before further ReadDir
// calls.
func (d *Dir) AddEntry(newEntry *DirEntry) (err error) {
	defer essentials.AddCtxTo("AddEntry", &err)
	if err := d.Reset(); err != nil {
		return err
	}
	for offset := int64(0); true; offset += 1 {
		cluster, err := d.chain.ReadCluster()
		if err != nil {
			return err
		}
		for i := 0; i < len(cluster); i += 32 {
			var entry DirEntry
			copy(entry[:], cluster[i:])
			if entry.IsFree() {
				copy(cluster[i:], newEntry[:])
				return d.chain.WriteCluster(cluster)
			}
		}
		if newOffset, err := d.chain.Seek(1, io.SeekCurrent); err != nil {
			return err
		} else if newOffset == offset {
			break
		}
	}

	// Create a new cluster of listings.
	if err := d.chain.Extend(); err != nil {
		return err
	}
	cluster := make([]byte, d.chain.fs.ClusterSize())
	copy(cluster, newEntry[:])
	return d.chain.WriteCluster(cluster)
}
