package fatfs

import (
	"testing"
	"time"
)

func TestRemove(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO")
	if err != nil {
		t.Fatal(err)
	}
	dir := NewDir(RootDirChain(fs))
	newDir, err := Mkdir(dir, "DIR", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	newDir, err = Mkdir(newDir, "SUBDIR", time.Now())
	if err != nil {
		t.Fatal(err)
	}

	cluster, err := fs.Alloc()
	if err != nil {
		t.Fatal(err)
	}
	newDir.AddEntry(NewDirEntry("FOO.TXT", cluster, 13, time.Now(), false))

	if err := Remove(dir, "DIR"); err != nil {
		t.Fatal(err)
	}

	listing, err := dir.ReadDir()
	if err != nil {
		t.Fatal(err)
	} else if len(listing) != 0 {
		t.Error("unexpected entries:", listing)
	}

	// Make sure that clusters 3 and onward are now free.
	for i := 3; i < 100; i++ {
		cluster, err := fs.Alloc()
		if err != nil {
			t.Fatal(err)
		}
		if cluster != uint32(i) {
			t.Errorf("expected %d but got %d", i, cluster)
		}
	}
}
