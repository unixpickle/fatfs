package main

import (
	"io/ioutil"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/fatfs"
)

func main() {
	dev := make(fatfs.RAMDisk, 299999744)
	fs, err := fatfs.FormatFS(dev, "TEST")
	essentials.Must(err)

	fileChunk := createFileCluster(fs)
	entry := createDirEntry(fileChunk)
	rootDir := fatfs.NewDir(fatfs.RootDirChain(fs))
	rootDir.AddEntry(entry)

	ioutil.WriteFile("disk.img", dev, 0755)
}

func createFileCluster(fs *fatfs.FS) uint32 {
	fileCluster, err := fs.Alloc()
	essentials.Must(err)
	fileChain := fatfs.NewChain(fs, fileCluster)
	fileContents := make([]byte, fs.ClusterSize())
	copy(fileContents, []byte("Hello, world!"))
	fileChain.WriteCluster(fileContents)
	return fileCluster
}

func createDirEntry(fileChunk uint32) *fatfs.DirEntry {
	var res fatfs.DirEntry
	copy(res.Name(), []byte("README  TXT"))
	res.SetFstClusLO(uint16(fileChunk))
	res.SetFstClusHI(uint16(fileChunk >> 16))
	res.SetFileSize(13)
	return &res
}
