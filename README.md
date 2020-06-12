# gin-timeout
Timeout Middleware for Gin framework

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
	engine.Use(timeout.Timeout(time.Second * 2))

	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a handler that will last 5 seconds
	engine.GET("/long", long)

	// create a handler that will last 5 seconds but can be canceled.
	engine.GET("/long2", long2)

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

func long2(c *gin.Context) {
	if doSomething(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{"hello": "long2"})
	}
}

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
Date: Thu, 11 Jun 2020 06:33:48 GMT
Content-Length: 45
Content-Type: text/plain; charset=utf-8

{"code":"E509","msg":"http: Handler timeout"}%
```
```
╰─$ curl -i http://localhost:8080/long2
HTTP/1.1 503 Service Unavailable
Date: Fri, 12 Jun 2020 12:31:14 GMT
Content-Length: 45
Content-Type: text/plain; charset=utf-8

{"code":"E509","msg":"http: Handler timeout"}
```

```
╰─$ curl -i http://localhost:8080/short
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 12 Jun 2020 12:32:13 GMT
Content-Length: 17

{"hello":"short"}
```