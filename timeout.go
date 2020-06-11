package timeout

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout/buffpool"
	"net/http"
	"time"
)

const (
	HandlerFuncTimeout = "E509"
	ErrUnknowError     = "E003"
)

func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// sync.Pool
		buffer := buffpool.GetBuff()

		tw := &TimeoutWriter{body: buffer, ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// Channel capacity must be greater than 0.
		// Otherwise, if the parent coroutine quit due to timeout,
		// the child coroutine may never be able to quit.
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			c.Abort()
			tw.ResponseWriter.WriteHeader(http.StatusInternalServerError)
			bt, _ := json.Marshal(errResponse{Code: ErrUnknowError,
				Msg: fmt.Sprintf("unknow internal error, %v", p)})
			tw.ResponseWriter.Write(bt)

		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
			bt, _ := json.Marshal(errResponse{Code: HandlerFuncTimeout,
				Msg: http.ErrHandlerTimeout.Error()})
			tw.ResponseWriter.Write(bt)
			c.Abort()
			tw.timedOut = true
			// If timeout happen, the buffer cannot be cleared actively,
			// but wait for the GC to recycle.
		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(buffer.Bytes())
			buffpool.PutBuff(buffer)
		}
	}
}

type errResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
