package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vearne/gin-timeout"
	"log"
	"time"
)

func main() {
	router := gin.Default()
	router.Use(timeout.Timeout(10*time.Second, `{"code": -1, "msg":"http: Handler timeout"}`))
	router.Static("/foo", "/tmp/foo")
	log.Fatal(router.Run(":8080"))
}

// mkdir -p /tmp/foo
// echo "a" >> /tmp/foo/a

// test case1:
// curl -I http://localhost:8080/foo/a
// test case2:
// curl -i http://localhost:8080/foo/a
