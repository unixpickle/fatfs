package fatfs

import (
	"testing"
)

func TestFormatFS(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	_, err := FormatFS(dev, "FOO", false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAlloc(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO", false)
	if err != nil {
		t.Fatal(err)
	}
	for i := uint32(3); i < fs.NumClusters(); i++ {
		clus, err := fs.Alloc()
		if err != nil {
			t.Fatal(err)
		}
		if i != clus {
			t.Fatalf("expected %d but got %d", i, clus)
		}
	}
	if _, err := fs.Alloc(); err == nil {
		t.Fatal("expected allocation failure")
	}
}
