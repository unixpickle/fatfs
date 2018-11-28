//+build !linux,!darwin

package fatfs

import (
	"io"
	"os"
)

func getDeviceSize(f *os.File) (int64, error) {
	return f.Seek(0, io.SeekEnd)
}
