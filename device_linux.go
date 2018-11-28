//+build linux
package fatfs

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

func getDeviceSize(f *os.File) (int64, error) {
	offset, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	if offset != 0 {
		return offset, nil
	}

	const IOCTL_BLKGETSIZE64 = 2148012658
	var size int64
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), IOCTL_BLKGETSIZE64,
		uintptr(unsafe.Pointer(&size)))
	if errno == 0 {
		return size, nil
	} else {
		return 0, errno
	}
}
