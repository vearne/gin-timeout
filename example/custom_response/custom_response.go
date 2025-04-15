package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	timeout "github.com/vearne/gin-timeout"
)

// test usage:
//  curl -H 'Accept-Language: zh' -i http://localhost:8080/short
//  curl -H 'Accept-Language: zh' -i http://localhost:8080/long
//  curl -H 'Accept-Language: zh' -i http://localhost:8080/long2
//  curl -H 'Accept-Language: zh' -i http://localhost:8080/long3

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("[start]AccessLog")
		ctx.Next()
		log.Println("[end]AccessLog")
	}
}

type Response struct {
	timeout.BaseResponse
}

// GetMessage rewrite method for other logic such as response translation
func (r *Response) GetContent(c *gin.Context) any {
	// You can translate the response message based on the request header
	log.Println("Accept-Language: ", c.Request.Header.Get("Accept-Language"))
	return r.Content
}

func main() {

	// create new gin without any middleware
	engine := gin.Default()

	defaultMsg := `{"code": -1, "msg":"http: Handler timeout with custom response"}`
	// add timeout middleware with 2 second duration
	engine.Use(timeout.Timeout(
		timeout.WithTimeout(2*time.Second),
		timeout.WithResponse(&Response{
			timeout.BaseResponse{
				Code:        http.StatusInternalServerError,
				Content:     defaultMsg,
				ContentType: "application/json; charset=utf-8",
			},
		}), // optional
		timeout.WithCallBack(func(r *http.Request) {
			fmt.Println("timeout happen, url:", r.URL.String())
		}), // optional
	))
	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a handler that will last 5 seconds
	engine.GET("/long", AccessLog(), long)

	// create a handler that will last 5 seconds but can be canceled.
	engine.GET("/long2", long2)

	// create a handler that will last 20 seconds but can be canceled.
	engine.GET("/long3", long3)

	engine.GET("/boundary", boundary)

	// run the server
	log.Fatal(engine.Run(":8080"))
}

func short(c *gin.Context) {
	defer func(writer gin.ResponseWriter) {
		fmt.Printf("c.Writer.Size: %v, %T\n", writer.Size(), writer)
	}(c.Writer)

	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "short"})
}

func long(c *gin.Context) {
	defer func(writer gin.ResponseWriter) {
		fmt.Printf("c.Writer.Size: %v, %T\n", writer.Size(), writer)
	}(c.Writer)

	fmt.Println("handler-long1, do something...")
	time.Sleep(3 * time.Second)
	fmt.Println("handler-long2, do something...")
	time.Sleep(3 * time.Second)
	fmt.Println("handler-long3, do something...")
	c.JSON(http.StatusOK, gin.H{"hello": "long"})
}

func boundary(c *gin.Context) {
	time.Sleep(2 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "boundary"})
}

func long2(c *gin.Context) {
	if doSomething(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{"hello": "long2"})
	}
}

func long3(c *gin.Context) {
	// request a slow service
	// see  https://github.com/vearne/gin-timeout/blob/master/example/slow_service.go
	url := "http://localhost:8882/hello"
	// Notice:
	// Please use c.Request.Context(), the handler will be canceled where timeout event happen.
	req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	client := http.Client{Timeout: 100 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Where timeout event happen, a error will be received.
		fmt.Println("error1:", err)
		return
	}
	defer resp.Body.Close()
	s, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error2:", err)
		return
	}
	fmt.Println(s)
}

// A cancelCtx can be canceled.
// When canceled, it also cancels any children that implement canceler.
func doSomething(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		fmt.Println("doSomething is canceled.")
		return false
	case <-time.After(5 * time.Second):
		fmt.Println("doSomething is done.")
		return true
	}
}
