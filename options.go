package timeout

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CallBackFunc func(*http.Request)
type GinCtxCallBackFunc func(*gin.Context)
type Option func(*TimeoutWriter)

type TimeoutOptions struct {
	CallBack       CallBackFunc
	GinCtxCallBack GinCtxCallBackFunc
	Timeout        time.Duration
	Response       Response
}

func WithTimeout(d time.Duration) Option {
	return func(t *TimeoutWriter) {
		t.Timeout = d
	}
}

// Optional parameters
func WithErrorHttpCode(code int) Option {
	return func(t *TimeoutWriter) {
		if t.Response == nil {
			t.Response = defaultResponse
		}
		t.Response.SetCode(code)
	}
}

// Optional parameters
func WithDefaultMsg(resp interface{}) Option {
	return func(t *TimeoutWriter) {
		if t.Response == nil {
			t.Response = defaultResponse
		}
		t.Response.SetContent(resp)
	}
}

// Optional parameters
func WithContentType(ct string) Option {
	return func(t *TimeoutWriter) {
		if t.Response == nil {
			t.Response = defaultResponse
		}
		t.Response.SetContentType(ct)
	}
}

func WithResponse(resp Response) Option {
	return func(t *TimeoutWriter) {
		if resp != nil {
			t.Response = resp
		}
	}
}

// Optional parameters
func WithCallBack(f CallBackFunc) Option {
	return func(t *TimeoutWriter) {
		t.CallBack = f
	}
}

// Optional parameters
func WithGinCtxCallBack(f GinCtxCallBackFunc) Option {
	return func(t *TimeoutWriter) {
		t.GinCtxCallBack = f
	}
}
