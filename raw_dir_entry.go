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

// NewRawDirEntry creates a RawDirEntry given some
// meta-data about a file.
//
// The name must be formatted with FormatName().
func NewRawDirEntry(name string, cluster, size uint32, creation time.Time, dir bool) *RawDirEntry {
	if len(name) != 11 {
		panic("invalid name argument")
	}
	var res RawDirEntry
	copy(res.Name(), []byte(name))
	res.SetFstClusLO(uint16(cluster))
	res.SetFstClusHI(uint16(cluster >> 16))
	res.SetFileSize(size)
	res.SetCrtDate(fatDate(creation))
	res.SetCrtTime(fatTime(creation))
	res.SetWrtDate(fatDate(creation))
	res.SetWrtTime(fatTime(creation))
	res.SetLstAccDate(fatDate(creation))
	if dir {
		res.SetAttr(Directory)
	}
	return &res
}

// FirstCluster gets the first cluster of the file.
func (r *RawDirEntry) FirstCluster() uint32 {
	return uint32(r.FstClusLO()) | (uint32(r.FstClusHI()) << 16)
}

// IsFree checks if the directory entry is a free slot.
func (r *RawDirEntry) IsFree() bool {
	return r.Name()[0] == 0 || r.Name()[0] == 0xe5
}

// IsLongName checks if this is a long name entry.
func (r *RawDirEntry) IsLongName() bool {
	return r.Attr()&0x3f == LongName
}

// IsDotPointer checks if this is "." or "..".
func (r *RawDirEntry) IsDotPointer() bool {
	return string(r.Name()) == ".          " || string(r.Name()) == "..         "
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

// UnformatName turns a filename like "FOO     TXT" into a
// normal name like "FOO.TXT".
func UnformatName(name string) string {
	if len(name) != 11 {
		panic("invalid name")
	}
	if name[8:] == "   " {
		return strings.TrimRight(name[:8], " ")
	} else {
		return strings.TrimRight(name[:8], " ") + "." + strings.TrimRight(name[8:], " ")
	}
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
