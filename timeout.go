package timeout

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout/buffpool"
	"net/http"
	"time"
)

func Timeout(t time.Duration, defaultMsg string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// sync.Pool
		buffer := buffpool.GetBuff()

		tw := &TimeoutWriter{body: buffer, ResponseWriter: c.Writer,
			h: make(http.Header), defaultMsg: defaultMsg}
		c.Writer = tw

		// Restore Context's writer
		defer func() {
			c.Writer = tw.ResponseWriter
		}()

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
			panic(p)

		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.timedOut = true
			tw.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
			tw.ResponseWriter.Write([]byte(tw.errorBody()))
			c.Abort()

			// If timeout happen, the buffer cannot be cleared actively,
			// but wait for the GC to recycle.
		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}

			if !tw.wroteHeader {
				tw.code = http.StatusOK
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(buffer.Bytes())
			buffpool.PutBuff(buffer)
		}

	}
}
