package timeout

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
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
	timedOut    bool
	wroteHeader bool
}

func (tw *TimeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, nil
	}

	return tw.body.Write(b)
}

func (tw *TimeoutWriter) WriteHeader(code int) {
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return
	}
	tw.writeHeader(code)
}

func (tw *TimeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *TimeoutWriter) WriteHeaderNow() {}

func (tw *TimeoutWriter) Header() http.Header {
	return tw.h
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
