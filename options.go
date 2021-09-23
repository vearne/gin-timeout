package timeout

import (
	"net/http"
	"time"
)

type CallBackFunc func(*http.Request)

type SkipFunc func(c *gin.Context) bool

type Option func(*TimeoutWriter)

type TimeoutOptions struct {
	CallBack      CallBackFunc
	DefaultMsg    string
	Timeout       time.Duration
	ErrorHttpCode int
	SkipFunc      SkipFunc
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
func WithDefaultMsg(s string) Option {
	return func(t *TimeoutWriter) {
		t.DefaultMsg = s
	}
}

// Optional parameters
func WithCallBack(f CallBackFunc) Option {
	return func(t *TimeoutWriter) {
		t.CallBack = f
	}
}

func WithSkip(f SkipFunc) Option {
	return func(t *TimeoutWriter) {
		t.SkipFunc = f
	}
}