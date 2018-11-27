package fatfs

import (
	"fmt"
	"testing"
	"time"
)

func TestAddEntry(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO")
	if err != nil {
		t.Fatal(err)
	}
	dir := NewDir(RootDirChain(fs))
	for i := 0; i < 384; i++ {
		contents, err := fs.Alloc()
		if err != nil {
			t.Fatal(err)
		}
		dir.AddEntry(NewDirEntry(fmt.Sprintf("%d.TXT", i), contents, uint32(i%15), time.Now(),
			false))
	}
	listings, err := dir.ReadDir()
	if err != nil {
		t.Fatal(err)
	}
	if len(listings) != 384 {
		t.Errorf("unexpected length: %d", len(listings))
	}
	for i, listing := range listings {
		name := fmt.Sprintf("%d.TXT", i)
		if string(listing.Name()) != name {
			t.Errorf("expected name %s but got %s", name, string(listing.Name()))
		}
	}
}

func TestRemoveEntry(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO")
	if err != nil {
		t.Fatal(err)
	}
	dir := NewDir(RootDirChain(fs))
	for i := 0; i < 384; i++ {
		contents, err := fs.Alloc()
		if err != nil {
			t.Fatal(err)
		}
		dir.AddEntry(NewDirEntry(fmt.Sprintf("%d.TXT", i), contents, uint32(i%15), time.Now(),
			false))
	}
	if _, err := dir.RemoveEntry("123.TXT"); err != nil {
		t.Fatal(err)
	}
	if _, err := dir.RemoveEntry("123.TXT"); err == nil {
		t.Fatal(err)
	}
	listings, err := dir.ReadDir()
	if err != nil {
		t.Fatal(err)
	}
	if len(listings) != 383 {
		t.Errorf("unexpected length: %d", len(listings))
	}
	for i, listing := range listings {
		if i >= 123 {
			i += 1
		}
		name := fmt.Sprintf("%d.TXT", i)
		if string(listing.Name()) != name {
			t.Errorf("expected name %s but got %s", name, string(listing.Name()))
		}
	}
}
