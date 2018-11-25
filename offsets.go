package fatfs

// ClusterSector gets the first sector of the cluster.
func ClusterSector(b *BootSector, cluster uint32) uint32 {
	firstData := uint32(b.RsvdSecCnt()) + uint32(b.NumFATs())*b.FatSz32()
	return firstData + (cluster-2)*uint32(b.SecPerClus())
}

// RootDirSector gets the first sector of the root
// directory.
func RootDirSector(b *BootSector) uint32 {
	return ClusterSector(b, b.RootClus())
}
