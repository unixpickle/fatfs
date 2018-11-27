package fatfs

import (
	"testing"
	"time"
)

func TestLongName(t *testing.T) {
	entry := NewDirEntry("this is a long name", 0, 13, time.Now(), false)
	if entry.Name() != "this is a long name" {
		t.Error("unexpected name:", entry.Name())
	}
}
