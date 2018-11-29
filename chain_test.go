package fatfs

import (
	"io"
	"testing"
)

func TestChainExtend(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO", false)
	if err != nil {
		t.Fatal(err)
	}
	chain := RootDirChain(fs)
	for i := 0; i < 4; i++ {
		if err := chain.Extend(); err != nil {
			t.Fatal(err)
		}
		verifyCluster(t, chain)
	}
	if offset, err := chain.Seek(0, io.SeekCurrent); err != nil {
		t.Fatal(err)
	} else if offset != 4 {
		t.Fatalf("expected offset 4 but got %d", offset)
	}
	verifyCluster(t, chain)
}

func TestChainSeek(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO", false)
	if err != nil {
		t.Fatal(err)
	}
	chain := RootDirChain(fs)
	for i := 0; i < 4; i++ {
		if err := chain.Extend(); err != nil {
			t.Fatal(err)
		}
	}
	verifyCluster(t, chain)
	if offset, err := chain.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if offset != 0 {
		t.Errorf("expected offset 0 but got %d", offset)
	}
	verifyCluster(t, chain)
	if offset, err := chain.Seek(10, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if offset != 4 {
		t.Errorf("expected offset 4 but got %d", offset)
	}
	verifyCluster(t, chain)
	if offset, err := chain.Seek(-3, io.SeekCurrent); err != nil {
		t.Fatal(err)
	} else if offset != 1 {
		t.Errorf("expected offset 1 but got %d", offset)
	}
	verifyCluster(t, chain)
	if offset, err := chain.Seek(2, io.SeekStart); err != nil {
		t.Fatal(err)
	} else if offset != 2 {
		t.Errorf("expected offset 2 but got %d", offset)
	}
	verifyCluster(t, chain)
	if offset, err := chain.Seek(-3, io.SeekEnd); err != nil {
		t.Fatal(err)
	} else if offset != 1 {
		t.Errorf("expected offset 1 but got %d", offset)
	}
	verifyCluster(t, chain)
}

func TestChainTrunc(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO", false)
	if err != nil {
		t.Fatal(err)
	}
	chain := RootDirChain(fs)
	for i := 0; i < 4; i++ {
		if err := chain.Extend(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := chain.Seek(1, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if err := chain.Truncate(); err != nil {
		t.Fatal(err)
	}
	verifyCluster(t, chain)
	if offset, err := chain.Seek(0, io.SeekCurrent); err != nil {
		t.Fatal(err)
	} else if offset != 3 {
		t.Errorf("expected offset 3 but got %d", offset)
	}
	verifyCluster(t, chain)
	if err := chain.Extend(); err != nil {
		t.Fatal(err)
	}
	verifyCluster(t, chain)
}

func TestChainFree(t *testing.T) {
	dev := make(RAMDisk, 4096*80000)
	fs, err := FormatFS(dev, "FOO", false)
	if err != nil {
		t.Fatal(err)
	}
	chain := RootDirChain(fs)
	for i := 0; i < 4; i++ {
		if err := chain.Extend(); err != nil {
			t.Fatal(err)
		}
		verifyCluster(t, chain)
	}
	if err := chain.Free(); err != nil {
		t.Fatal(err)
	}
	cluster, err := fs.Alloc()
	if err != nil {
		t.Fatal(err)
	}
	chain = NewChain(fs, cluster)
	verifyCluster(t, chain)
	for i := 0; i < 4; i++ {
		if err := chain.Extend(); err != nil {
			t.Fatal(err)
		}
		verifyCluster(t, chain)
	}
}

func verifyCluster(t *testing.T, c *Chain) {
	expected := uint32(len(c.prev) + 2)
	if c.cluster != expected {
		t.Errorf("expected cluster %d but got %d", expected, c.cluster)
	}
}
