package fatfs

import (
	"time"

	"github.com/unixpickle/essentials"
)

// Mkdir creates an empty directory.
func Mkdir(parent *Dir, name string, date time.Time) (d *Dir, err error) {
	defer essentials.AddCtxTo("Mkdir", &err)

	fs := parent.Chain.FS()
	dirCluster, err := fs.Alloc()
	if err != nil {
		return nil, err
	}
	chain := NewChain(fs, dirCluster)
	clusterData := make([]byte, fs.ClusterSize())
	copy(clusterData, NewRawDirEntry(".          ", dirCluster, 0, date, true)[:])
	copy(clusterData[32:], NewRawDirEntry("..         ", parent.Chain.FirstCluster(), 0, date,
		true)[:])
	if err := chain.WriteCluster(clusterData); err != nil {
		return nil, err
	}

	entry := NewDirEntry(name, dirCluster, 0, date, true)
	if err := parent.AddEntry(entry); err != nil {
		fs.WriteFAT(dirCluster, 0)
		return nil, err
	}

	return NewDir(chain), nil
}

// Remove deletes a file or directory.
// It uses recursion if necessary.
func Remove(parent *Dir, name string) (err error) {
	defer essentials.AddCtxTo("Remove", &err)
	entry, err := parent.RemoveEntry(name)
	if err != nil {
		return err
	}
	chain := NewChain(parent.Chain.FS(), entry.Raw().FirstCluster())
	if entry.Raw().Attr()&Directory == Directory {
		dir := NewDir(chain)
		listing, err := dir.ReadDir()
		if err != nil {
			return err
		}
		for _, entry := range listing {
			if !entry.Raw().IsDotPointer() {
				if err := Remove(dir, entry.Name()); err != nil {
					return err
				}
			}
		}
	}
	return chain.Free()
}
