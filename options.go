package timeout

import (
	"net/http"
	"time"
)

type CallBackFunc func(*http.Request)
type Option func(*TimeoutWriter)

type TimeoutOptions struct {
	CallBack      CallBackFunc
	DefaultMsg    interface{}
	Timeout       time.Duration
	ErrorHttpCode int
}

func WithTimeout(d time.Duration) Option {
	return func(t *TimeoutWriter) {
		t.Timeout = d
	}
}

// Optional parameters
func WithErrorHttpCode(code int) Option {
	return func(t *TimeoutWriter) {
		t.ErrorHttpCode = code
	}
}

// Optional parameters
func WithDefaultMsg(resp interface{}) Option {
	return func(t *TimeoutWriter) {
		t.DefaultMsg = resp
	}
}

// Optional parameters
func WithCallBack(f CallBackFunc) Option {
	return func(t *TimeoutWriter) {
		t.CallBack = f
	}
}
