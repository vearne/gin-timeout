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

	// create a route that will last 5 seconds
	engine.GET("/long", long)

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
╰─$ curl -i http://localhost:8080/short
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 11 Jun 2020 06:34:13 GMT
Content-Length: 17

{"hello":"short"}
```

