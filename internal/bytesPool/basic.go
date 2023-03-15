package bytesPool

import (
	"bytes"
	"sync"
)

var Pool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}
