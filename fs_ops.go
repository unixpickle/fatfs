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
	clusterData := make([]byte, fs.ClusterSize())
	copy(clusterData, NewDirEntry(".          ", dirCluster, 0, date, true)[:])
	copy(clusterData[32:], NewDirEntry("..         ", parent.Chain.FirstCluster(), 0, date,
		true)[:])

	entry := NewDirEntry(FormatName(name), dirCluster, 0, date, true)
	if err := parent.AddEntry(entry); err != nil {
		fs.WriteFAT(dirCluster, 0)
		return nil, err
	}

	return NewDir(NewChain(fs, dirCluster)), nil
}
