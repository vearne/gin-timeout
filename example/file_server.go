package main

import (
	"github.com/gin-gonic/gin"
	timeout "github.com/vearne/gin-timeout"
	"log"
	"time"
)

func main() {
	router := gin.Default()
	defaultMsg := `{"code": -1, "msg":"http: Handler timeout"}`
	router.Use(timeout.Timeout(timeout.WithTimeout(10*time.Second),
		timeout.WithDefaultMsg(defaultMsg)))
	//router.StaticFS("/static", gin.Dir("/tmp/static", true))
	router.Static("/static", "/tmp/static")
	log.Fatal(router.Run(":8080"))
}

// mkdir -p /tmp/foo
// echo "a" >> /tmp/foo/a

// test case1:
// curl -I http://localhost:8080/static/a
// test case2:
// curl -i http://localhost:8080/static/a
