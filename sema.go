package fastchan

import (
	_ "sync"
	_ "unsafe"
)

// Trick to take unexported functions from sync package, combined with empty assembly to allow function declaration
// without function body

//go:linkname semacquire sync.runtime_Semacquire
func semacquire(addr *uint32)

//go:linkname semrelease sync.runtime_Semrelease
func semrelease(addr *uint32, handoff bool, skipframes int)
