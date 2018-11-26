package main

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/fatfs"
)

func main() {
	dev := make(fatfs.RAMDisk, 299999744)
	fs, err := fatfs.FormatFS(dev, "TEST")
	essentials.Must(err)

	fileCluster, fileSize := createFileCluster(fs)
	rootDir := fatfs.NewDir(fatfs.RootDirChain(fs))
	subDir, err := fatfs.Mkdir(rootDir, "photos", time.Now())
	essentials.Must(err)
	subDir.AddEntry(fatfs.NewDirEntry("example.jpg", fileCluster, fileSize, time.Now(), false))

	ioutil.WriteFile("disk.img", dev, 0755)
}

func createFileCluster(fs *fatfs.FS) (uint32, uint32) {
	fileCluster, err := fs.Alloc()
	essentials.Must(err)
	fileChain := fatfs.NewChain(fs, fileCluster)

	inFile, err := os.Open("example.jpg")
	essentials.Must(err)
	defer inFile.Close()

	size, err := fileChain.ReadFrom(inFile)
	essentials.Must(err)

	return fileCluster, uint32(size)
}
