package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"time"
)

func main() {
	r := gin.Default()
	//  This handler will take 20 seconds
	r.GET("/hello", func(c *gin.Context) {
		beginTime := time.Now()
		counter := 20
		// 100 byte
		c.Status(200)
		c.Header("content-length", strconv.Itoa(counter))
		c.Header("content-type", "text/plain")
		for counter > 0 {
			time.Sleep(1 * time.Second)
			log.Println(strconv.Itoa(counter))
			c.Writer.WriteString(strconv.Itoa(counter))
			counter--
		}
		log.Println("handler cost:", time.Since(beginTime))
	})
	r.Run(":8882")
}
