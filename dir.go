package fatfs

// RootDirSector gets the sector of the root directory.
func RootDirSector(b *BootSector) uint32 {
	return uint32(b.RsvdSecCnt()) + uint32(b.NumFATs())*b.FatSz32() + b.RootClus()
}
