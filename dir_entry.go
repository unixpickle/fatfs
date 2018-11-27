package fatfs

import (
	"time"

	"github.com/unixpickle/essentials"
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
	checksum := shortNameChecksum(short.Name())
	words := runesToWords([]rune(name))
	numParts := len(words) / 13
	if len(words)%13 != 0 {
		numParts += 1
	}
	parts := DirEntry{}
	for i := 0; i < numParts; i++ {
		last := i == numParts-1
		endIdx := i*13 + 13
		if endIdx >= len(words) {
			endIdx = len(words)
		}
		parts = append(parts, packLongEntry(words[i*13:endIdx], i+1, last, checksum))
	}
	essentials.Reverse(parts)
	return append(parts, short)
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
	var words []uint16
	for i := len(d) - 2; i >= 0; i-- {
		words = append(words, unpackLongEntry(d[i])...)
	}
	return string(wordsToRunes(words))
}

func unpackLongEntry(raw *RawDirEntry) []uint16 {
	var res []uint16
	for _, byteRange := range [][2]int{{1, 11}, {14, 26}, {28, 32}} {
		for i := byteRange[0]; i < byteRange[1]; i += 2 {
			word := Endian.Uint16(raw[i : i+2])
			if word == 0 {
				return res
			}
			res = append(res, word)
		}
	}
	return res
}

func packLongEntry(data []uint16, idx int, last bool, checksum uint8) *RawDirEntry {
	var res RawDirEntry
	res.SetAttr(LongName)
	res[0] = uint8(idx)
	if last {
		res[0] |= 0x40
	}
	res[13] = checksum
	var dataIdx int
	for _, byteRange := range [][2]int{{1, 11}, {14, 26}, {28, 32}} {
		for i := byteRange[0]; i < byteRange[1]; i += 2 {
			var word uint16
			if dataIdx < len(data) {
				word = data[dataIdx]
			} else if dataIdx == len(data) {
				word = 0
			} else {
				word = 0xffff
			}
			Endian.PutUint16(res[i:i+2], word)
			dataIdx++
		}
	}
	return &res
}

func wordsToRunes(words []uint16) []rune {
	var res []rune
	for _, word := range words {
		res = append(res, rune(word))
	}
	return res
}

func runesToWords(runes []rune) []uint16 {
	var res []uint16
	for _, r := range runes {
		res = append(res, uint16(r))
	}
	return res
}

func shortNameChecksum(name []byte) uint8 {
	var res uint8
	for _, b := range name {
		var addition uint8
		if res&1 != 0 {
			addition = 0x80
		}
		res = addition + (res >> 1) + b
	}
	return res
}
