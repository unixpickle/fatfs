package main

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/fatfs"
)

func main() {
	dev := readImage()
	fs, err := fatfs.NewFS(dev)
	essentials.Must(err)

	os.RemoveAll("extracted")
	essentials.Must(os.Mkdir("extracted", 0755))

	rootDir := fatfs.NewDir(fatfs.RootDirChain(fs))
	extractDirectory("extracted", rootDir)
}

func readImage() fatfs.RAMDisk {
	file, err := os.Open("disk.img.gz")
	essentials.Must(err)
	defer file.Close()
	reader, err := gzip.NewReader(file)
	essentials.Must(err)
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)
	essentials.Must(err)
	return data
}

func extractDirectory(dest string, source *fatfs.Dir) {
	listing, err := source.ReadDir()
	essentials.Must(err)
	for _, entry := range listing {
		if string(entry.Name()) == ".          " || string(entry.Name()) == "..         " {
			continue
		}
		outName := fatfs.UnformatName(string(entry.Name()))
		chain := fatfs.NewChain(source.Chain.FS(), entry.FirstCluster())
		filePath := filepath.Join(dest, outName)
		if entry.Attr()&fatfs.Directory == fatfs.Directory {
			essentials.Must(os.Mkdir(filePath, 0755))
			extractDirectory(filePath, fatfs.NewDir(chain))
		} else {
			f, err := os.Create(filePath)
			essentials.Must(err)
			_, err = chain.WriteTo(f)
			essentials.Must(err)
			f.Truncate(int64(entry.FileSize()))
			f.Close()
		}
	}
}
