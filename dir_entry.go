package fatfs

import (
	"time"
)

// A DirEntry is a directory entry that potentially has a
// long name.
type DirEntry []*RawDirEntry

// NewDirEntry creates a DirEntry using a long name and
// a pre-existing raw entry.
func NewDirEntry(name string, cluster, size uint32, date time.Time, dir bool) DirEntry {
	return WrapDirEntry(name, NewRawDirEntry(FormatName(name), cluster, size, date, dir))
}

// WrapDirEntry creates a DirEntry around a RawDirEntry.
func WrapDirEntry(name string, short *RawDirEntry) DirEntry {
	if name == UnformatName(string(short.Name())) {
		return DirEntry{short}
	}
	// TODO: create a long entry.
	panic("not yet implemented")
}

// Raw gets the short entry corresponding to this entry.
// This can be used for all attributes besides the name.
func (d DirEntry) Raw() *RawDirEntry {
	return d[len(d)-1]
}

// Name gets the name of the directory entry. This may be
// the short name if no long name is present.
func (d DirEntry) Name() string {
	if len(d) == 1 {
		return UnformatName(string(d[0].Name()))
	}
	var runes []rune
	for i := len(d) - 2; i >= 0; i-- {
		runes = append(runes, unpackLongEntry(d[i])...)
	}
	return string(runes)
}

func unpackLongEntry(raw *RawDirEntry) []rune {
	var res []rune
	for _, byteRange := range [][2]int{{1, 11}, {14, 26}, {28, 32}} {
		for i := byteRange[0]; i < byteRange[1]; i += 2 {
			word := Endian.Uint16(raw[i : i+2])
			if word == 0 {
				return res
			}
			res = append(res, rune(word))
		}
	}
	return res
}
