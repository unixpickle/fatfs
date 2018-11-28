//+build darwin

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

	const IOCTL_DKIOCGETBLOCKCOUNT = 1074291737
	const IOCTL_DKIOCGETBLOCKSIZE = 1074029592
	var blockSize int64
	var blockCount int64
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), IOCTL_DKIOCGETBLOCKSIZE,
		uintptr(unsafe.Pointer(&blockSize)))
	if errno != 0 {
		return 0, errno
	}
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), IOCTL_DKIOCGETBLOCKCOUNT,
		uintptr(unsafe.Pointer(&blockCount)))
	if errno == 0 {
		return blockSize * blockCount, nil
	} else {
		return 0, errno
	}
}
