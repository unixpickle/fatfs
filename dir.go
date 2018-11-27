package fatfs

import (
	"errors"
	"io"
	"os"

	"github.com/unixpickle/essentials"
)

// A Dir is an open handle to a directory.
type Dir struct {
	Chain *Chain
}

// NewDir creates a Dir from a Chain.
// The Chain must be seeked to the start.
func NewDir(c *Chain) *Dir {
	return &Dir{Chain: c}
}

// ReadDirRaw reads the raw directory listings.
func (d *Dir) ReadDirRaw() (entries []*RawDirEntry, err error) {
	defer essentials.AddCtxTo("ReadDirRaw", &err)
	if _, err := d.Chain.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	for {
		cluster, done, err := d.Chain.ReadNext()
		if err != nil {
			return entries, err
		}
		for i := 0; i < len(cluster); i += 32 {
			var entry RawDirEntry
			copy(entry[:], cluster[i:])
			if !entry.IsFree() {
				entries = append(entries, &entry)
			}
		}
		if done {
			break
		}
	}
	return
}

// ReadDir reads the directory's entries.
func (d *Dir) ReadDir() (entries []DirEntry, err error) {
	defer essentials.AddCtxTo("ReadDir", &err)

	rawEntries, err := d.ReadDirRaw()
	if err != nil {
		return nil, err
	}

	var longEntry DirEntry
	for _, entry := range rawEntries {
		longEntry = append(longEntry, entry)
		if !entry.IsLongName() {
			entries = append(entries, longEntry)
			longEntry = DirEntry{}
		}
	}
	if len(longEntry) > 0 {
		return entries, errors.New("missing final short entry")
	}
	return entries, nil
}

// WriteDir updates the directory's entries.
func (d *Dir) WriteDir(entries []DirEntry) (err error) {
	defer essentials.AddCtxTo("WriteDir", &err)

	clusterSize := d.Chain.FS().ClusterSize()

	var clusters [][]byte
	var currentCluster []byte

	finishCurrent := func() {
		currentCluster = append(currentCluster, make([]byte, clusterSize-len(currentCluster))...)
		clusters = append(clusters, currentCluster)
		currentCluster = nil
	}

	for _, entry := range entries {
		var encoded []byte
		for _, rawEntry := range entry {
			encoded = append(encoded, rawEntry[:]...)
		}
		if len(encoded) > clusterSize {
			return errors.New("entry is too large to fit into a cluster")
		}
		if len(encoded)+len(currentCluster) > clusterSize {
			finishCurrent()
		}
		currentCluster = append(currentCluster, encoded...)
	}
	finishCurrent()

	return d.Chain.SetClusters(clusters)
}

// AddEntry adds a directory entry.
func (d *Dir) AddEntry(newEntry DirEntry) (err error) {
	defer essentials.AddCtxTo("AddEntry", &err)

	entries, err := d.ReadDir()
	if err != nil {
		return err
	}
	entries = append(entries, newEntry)
	return d.WriteDir(entries)
}

// RemoveEntry deletes the entry for the given name.
//
// If the entry is found, it is returned.
// If no entry is found, ErrNotExist is returned.
//
// The name is formatted; it is not a short name.
// Use "FOO.TXT", not "FOO     TXT".
func (d *Dir) RemoveEntry(name string) (entry DirEntry, err error) {
	entries, err := d.ReadDir()
	if err != nil {
		return nil, essentials.AddCtx("RemoveEntry", err)
	}
	for i, entry := range entries {
		if entry.Name() == name {
			essentials.OrderedDelete(&entries, i)
			return entry, essentials.AddCtx("RemoveEntry", d.WriteDir(entries))
		}
	}
	return nil, os.ErrNotExist
}

// AddRawEntry adds a raw directory entry.
func (d *Dir) AddRawEntry(entry *RawDirEntry) error {
	return d.AddEntry(DirEntry{entry})
}
