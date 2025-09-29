package timeout

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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

// WithErrorHttpCode Optional parameters
func WithErrorHttpCode(code int) Option {
	return func(t *TimeoutWriter) {
		if t.Response == nil {
			t.Response = defaultResponse
		}
		t.Response.SetCode(code)
	}
}

// WithDefaultMsg Optional parameters
func WithDefaultMsg(resp interface{}) Option {
	return func(t *TimeoutWriter) {
		if t.Response == nil {
			t.Response = defaultResponse
		}
		t.Response.SetContent(resp)
	}
}

// WithContentType Optional parameters
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

// WithCallBack Optional parameters
func WithCallBack(f CallBackFunc) Option {
	return func(t *TimeoutWriter) {
		t.CallBack = f
	}
}

// WithGinCtxCallBack Optional parameters
func WithGinCtxCallBack(f GinCtxCallBackFunc) Option {
	return func(t *TimeoutWriter) {
		t.GinCtxCallBack = f
	}
}
