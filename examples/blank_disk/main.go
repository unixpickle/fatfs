package main

import (
	"io/ioutil"

	"github.com/unixpickle/essentials"
	"github.com/unixpickle/fatfs"
)

func main() {
	dev := make(fatfs.RAMDisk, 299999744)
	_, err := fatfs.FormatFS(dev, "TEST")
	essentials.Must(err)
	ioutil.WriteFile("disk.img", dev, 0755)
}
