package reentrant

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

// GetCurrentGoroutineID returns the current goroutine id.
func GetCurrentGoroutineID() int64 {
	buf := make([]byte, 24)
	runtime.Stack(buf, false)

	i := bytes.IndexByte(buf[10:], ' ')
	b := buf[10 : 10+i]

	n, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("error occurred while getting goroutine id: %v", err))
	}
	return n
}
