# gin-timeout
![build_ci](https://travis-ci.com/vearne/gin-timeout.svg?branch=master)

针对gin的超时中间件


* [English README](https://github.com/vearne/gin-timeout/blob/master/README.md)

### 感谢
本库的实现受到了标准库
[http.TimeoutHandler](https://github.com/golang/go/blob/5f3dabbb79fb3dc8eea9a5050557e9241793dce3/src/net/http/server.go#L3255) 实现的启发

### 安装&使用
```
export GO111MODULE=on
go get github.com/vearne/gin-timeout
```

### 注意:
如果handler支持取消操作，那么需要传入context.Context为c.Request.Context()

### 示例
```
package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout"
	"io/ioutil"
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

	// create a handler that will last 20 seconds but can be canceled.
	engine.GET("/long3", long3)

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
	// 注意: 这里需要使用 c.Request.Context()
	// 当超时发生时，handler会被取消掉
	if doSomething(c.Request.Context()) {
		c.JSON(http.StatusOK, gin.H{"hello": "long2"})
	}
}

func long3(c *gin.Context) {
	// 我已经提供了一个慢服务的实例
	// 参见  https://github.com/vearne/gin-timeout/blob/master/example/slow_service.go
	url := "http://localhost:8882/hello"
	// 注意: 这里需要使用 c.Request.Context()
	// 当超时发生时，handler会被取消掉
	req, _ := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, url, nil)
	client := http.Client{Timeout: 100* time.Second}
	resp, err :=client.Do(req)
	if err != nil {
	    // 当超时发生时，会收到一个错误
		fmt.Println("error1:", err)
		return
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
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
```


