package fatfs

import (
	"path"
	"strings"
	"time"
)

// Directory attribute flags.
const (
	ReadOnly  = 0x01
	Hidden    = 0x02
	System    = 0x04
	VolumeID  = 0x08
	Directory = 0x10
	Archive   = 0x20
	LongName  = 0x0f
)

// NewDirEntry creates a DirEntry given some meta-data
// about a file.
func NewDirEntry(name string, cluster, size uint32, creation time.Time, isDir bool) *DirEntry {
	var res DirEntry
	copy(res.Name(), []byte(FormatName(name)))
	res.SetFstClusLO(uint16(cluster))
	res.SetFstClusHI(uint16(cluster >> 16))
	res.SetFileSize(size)
	res.SetCrtDate(fatDate(creation))
	res.SetCrtTime(fatTime(creation))
	res.SetWrtDate(fatDate(creation))
	res.SetWrtTime(fatTime(creation))
	res.SetLstAccDate(fatDate(creation))
	if isDir {
		res.SetAttr(Directory)
	}
	return &res
}

// IsFree checks if the directory entry is a free slot.
func (d *DirEntry) IsFree() bool {
	return d.Name()[0] == 0 || d.Name()[0] == 0xe5
}

// IsLongName checks if this is a long name entry.
func (d *DirEntry) IsLongName() bool {
	return d.Attr()&0x3f == LongName
}

// FormatName turns a regular filename, like "foo.txt",
// into a string like "FOO     TXT".
func FormatName(name string) string {
	name = strings.ToUpper(name)
	prefix := name
	ext := path.Ext(name)
	if ext != "" {
		prefix = name[:len(name)-len(ext)]
		ext = ext[1:]
	}
	prefix = spacePad(prefix, 8)
	ext = spacePad(ext, 3)
	return prefix + ext
}

func spacePad(str string, length int) string {
	if len(str) > length {
		return str[:length]
	}
	for len(str) < length {
		str += " "
	}
	return str
}

func fatDate(t time.Time) uint16 {
	return uint16(t.Day()) | (uint16(t.Month()) << 5) | ((uint16(t.Year()) - 1980) << 9)
}

func fatTime(t time.Time) uint16 {
	return (uint16(t.Second()) / 2) | (uint16(t.Minute()) << 5) | (uint16(t.Hour()) << 11)
}
