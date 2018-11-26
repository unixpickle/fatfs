package fatfs

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

// IsFree checks if the directory entry is a free slot.
func (d *DirEntry) IsFree() bool {
	return d.Name()[0] == 0 || d.Name()[0] == 0xe5
}

// IsLongName checks if this is a long name entry.
func (d *DirEntry) IsLongName() bool {
	return d.Attr()&0x3f == LongName
}
