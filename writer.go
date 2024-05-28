package timeout

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"sync/atomic"
)

type TimeoutWriter struct {
	gin.ResponseWriter
	// header
	h http.Header
	// body
	body           *bytes.Buffer
	TimeoutOptions // TimeoutOptions in options.go

	code        int
	mu          sync.Mutex
	timedOut    atomic.Bool
	wroteHeader atomic.Bool
	size        int
}

func (tw *TimeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut.Load() {
		return 0, nil
	}
	tw.size += len(b)
	return tw.body.Write(b)
}

func (tw *TimeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut.Load() {
		return
	}
	tw.writeHeader(code)
}

func (tw *TimeoutWriter) writeHeader(code int) {
	tw.wroteHeader.Store(true)
	tw.code = code
}

func (tw *TimeoutWriter) WriteHeaderNow() {}

func (tw *TimeoutWriter) Header() http.Header {
	return tw.h
}

func (tw *TimeoutWriter) Size() int {
	return tw.size
}
