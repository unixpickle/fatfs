package fatfs

// A DirEntry is a directory entry that potentially has a
// long name.
type DirEntry []*RawDirEntry

// NewDirEntry creates a DirEntry using a long name and
// a pre-existing raw entry.
func NewDirEntry(short *RawDirEntry, name string) DirEntry {
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
	// TODO: assemble the name here.
	panic("not yet implemented")
}
