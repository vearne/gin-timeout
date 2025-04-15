package timeout

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var defaultResponse = &BaseResponse{
	Code:        http.StatusServiceUnavailable,
	Content:     `{"code": -1, "msg":"http: Handler timeout"}`,
	ContentType: "text/plain; charset=utf-8",
}

type Response interface {
	GetCode(c *gin.Context) int
	GetContent(c *gin.Context) any
	GetContentType(c *gin.Context) string
	SetCode(int)
	SetContent(any)
	SetContentType(string)
}

type BaseResponse struct {
	Code        int
	Content     any
	ContentType string
}

func (r *BaseResponse) GetCode(c *gin.Context) int {
	return r.Code
}

func (r *BaseResponse) GetContent(c *gin.Context) any {
	return r.Content
}

func (r *BaseResponse) GetContentType(c *gin.Context) string {
	return r.ContentType
}

func (r *BaseResponse) SetCode(code int) {
	r.Code = code
}

func (r *BaseResponse) SetContent(content any) {
	r.Content = content
}

func (r *BaseResponse) SetContentType(contentType string) {
	r.ContentType = contentType
}
