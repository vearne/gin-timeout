package timeout

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout/buffpool"
	"net/http"
	"time"
)

var (
	defaultOptions TimeoutOptions
)

func init() {
	defaultOptions = TimeoutOptions{
		CallBack:      nil,
		DefaultMsg:    `{"code": -1, "msg":"http: Handler timeout"}`,
		Timeout:       3 * time.Second,
		ErrorHttpCode: http.StatusServiceUnavailable,
	}
}

func Timeout(opts ...Option) gin.HandlerFunc {
	return func(c *gin.Context) {
		// sync.Pool
		buffer := buffpool.GetBuff()

		tw := &TimeoutWriter{body: buffer, ResponseWriter: c.Writer,
			h: make(http.Header)}
		tw.TimeoutOptions = defaultOptions

		// Loop through each option
		for _, opt := range opts {
			// Call the option giving the instantiated
			opt(tw)
		}

		if tw.SkipFunc != nil {
			if tw.SkipFunc(c) {
				c.Next()
				return
			}
		}

		c.Writer = tw

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), tw.Timeout)
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

		var err error
		select {
		case p := <-panicChan:
			panic(p)

		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.timedOut = true
			tw.ResponseWriter.WriteHeader(tw.ErrorHttpCode)
			_, err = tw.ResponseWriter.Write([]byte(tw.DefaultMsg))
			if err != nil {
				panic(err)
			}
			c.Abort()

			// execute callback func
			if tw.CallBack != nil {
				tw.CallBack(c.Request.Clone(context.Background()))
			}
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
			_, err = tw.ResponseWriter.Write(buffer.Bytes())
			if err != nil {
				panic(err)
			}
			buffpool.PutBuff(buffer)
		}

	}
}
