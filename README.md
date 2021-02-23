# gin-timeout
Timeout Middleware for Gin framework

### Thanks
Inspired by golang source code [http.TimeoutHandler](https://github.com/golang/go/blob/5f3dabbb79fb3dc8eea9a5050557e9241793dce3/src/net/http/server.go#L3255)

### Usage
Download and install using go module:
```
export GO111MODULE=on
go get github.com/vearne/gin-timeout
```

### Example
```
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout"
	"log"
	"net/http"
	"time"
)

func main() {
	// create new gin without any middleware
	engine := gin.Default()

	// add timeout middleware with 2 second duration
	engine.Use(timeout.Timeout(time.Second*2, `{"code": -1, "msg":"http: Handler timeout"}`))

	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a handler that will last 5 seconds
	engine.GET("/long", long)

	// create a handler that will last 5 seconds but can be canceled.
	engine.GET("/long2", long2)

	engine.GET("/boundary", boundary)

	// run the server
	log.Fatal(engine.Run(":8080"))
}

func short(c *gin.Context) {
	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "short"})
}

func long(c *gin.Context) {
	time.Sleep(3 * time.Second)
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
```

### Output 
```
╰─$ curl -i http://localhost:8080/long
HTTP/1.1 503 Service Unavailable
Date: Tue, 23 Feb 2021 04:56:18 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

{"code": -1, "msg":"http: Handler timeout"}
```
```
╰─$ curl -i http://localhost:8080/long2
HTTP/1.1 503 Service Unavailable
Date: Tue, 23 Feb 2021 04:56:39 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

{"code": -1, "msg":"http: Handler timeout"}
```

```
╰─$ curl -i http://localhost:8080/short
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 23 Feb 2021 04:56:58 GMT
Content-Length: 17

{"hello":"short"}
```