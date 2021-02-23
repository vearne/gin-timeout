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
	// body
	body *bytes.Buffer
	// header
	h http.Header

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	code        int

	defaultMsg string
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

func (tw *TimeoutWriter) errorBody() string {
	if tw.defaultMsg != "" {
		return tw.defaultMsg
	}
	return "<html><head><title>Timeout</title></head><body><h1>Timeout</h1></body></html>"
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
