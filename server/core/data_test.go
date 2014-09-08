package core

import (
	"runtime"
	"testing"
)

func TestMemoryUse(t *testing.T) {
	SetDataPrefix(".")
	arr := make([]byte, 1<<20) // 1 MB
	t.Logf("Made byte array of length %d.\n", len(arr))
	ds := NewDataStore("testStore")
	name := "1mb"
	ds.Save(name, arr)
	arr, _ = ds.Load(name)
	startAlloc, startTotalAlloc := runtime.MemStats.Alloc, runtime.MemStats.TotalAlloc
	t.Logf("Mem stats before reloading and resaving:\n\tAlloc: %d\n\tTotalAlloc: %d\n", startAlloc, startTotalAlloc)
	for i := 0; i < 10; i++ { // if leaking, should leak at least 10 times the file size
		ds.Save(name, arr)
		arr, _ = ds.Load(name)
		runtime.GC()
	}
	endAlloc, endTotalAlloc := runtime.MemStats.Alloc, runtime.MemStats.TotalAlloc
	t.Logf("Mem stats after reloading and resaving:\n\tAlloc: %d\n\tTotalAlloc: %d\n", endAlloc, endTotalAlloc)
	t.Logf("Arr size: %d\n", len(arr))
	if endAlloc-startAlloc > uint64(len(arr)) || endTotalAlloc-startTotalAlloc > uint64(len(arr)) {
		t.Error("FAIL: Ending memory greatly exceeds starting memory")
		t.Fail()
	} else {
		t.Log("SUCCESS: No major memory leak in Data.")
	}
}
