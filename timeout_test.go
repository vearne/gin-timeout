package timeout

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func testEngine() *gin.Engine {
	engine := gin.Default()

	defaultMsg := `{"code": -1, "msg":"http: Handler timeout"}`
	// add timeout middleware with 2 second duration
	engine.Use(Timeout(
		WithTimeout(2*time.Second),
		WithErrorHttpCode(http.StatusRequestTimeout), // optional
		WithDefaultMsg(defaultMsg),                   // optional
		WithCallBack(func(r *http.Request) {
			fmt.Println("timeout happen, url:", r.URL.String())
		}), // optional
	))

	// create a handler that will last 1 seconds
	engine.GET("/short", short)

	// create a handler that will last 5 seconds
	engine.GET("/long", AccessLog(), long)

	engine.GET("/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, "this is page A")
	})
	engine.GET("/b", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/a")
	})

	return engine
}

func short(c *gin.Context) {
	defer func(writer gin.ResponseWriter) {
		fmt.Printf("c.Writer.Size: %v, %T\n", writer.Size(), writer)
	}(c.Writer)

	time.Sleep(1 * time.Second)
	c.JSON(http.StatusOK, gin.H{"hello": "short"})
}

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("[start]AccessLog")
		ctx.Next()
		log.Println("[end]AccessLog")
	}
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

func Get(uri string, router *gin.Engine, headers, querys map[string]string) (int, []byte) {
	u, _ := url.Parse(uri)
	q := u.Query()
	for k, v := range querys {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req := httptest.NewRequest("GET", u.String(), nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	body, _ := io.ReadAll(result.Body)
	return result.StatusCode, body
}

func TestTimeout(t *testing.T) {
	router := testEngine()
	var code int
	var b []byte
	code, _ = Get("/short", router, nil, nil)
	assert.Equal(t, http.StatusOK, code)

	code, b = Get("/long", router, nil, nil)
	assert.Equal(t, http.StatusRequestTimeout, code)
	assert.Equal(t, `{"code": -1, "msg":"http: Handler timeout"}`, string(b))

	code, _ = Get("/b", router, nil, nil)
	assert.Equal(t, http.StatusMovedPermanently, code)
}

func TestPanic(t *testing.T) {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		defer func() {
			if p := recover(); p != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					map[string]any{"code": -1, "msg": fmt.Sprintf("unknow error:%v", p)})
				return
			}
		}()
		c.Next()
	})
	router.Use(Timeout(WithTimeout(3 * time.Second)))
	router.GET("/panic", func(c *gin.Context) {
		time.Sleep(1 * time.Second)
		x := 0
		fmt.Println(100 / x)
	})

	code, b := Get("/panic", router, nil, nil)
	assert.Equal(t, http.StatusInternalServerError, code)
	assert.Contains(t, string(b), "integer divide by zero")
}

func TestWriteSize(t *testing.T) {
	router := gin.Default()
	router.Use(Timeout(WithTimeout(3 * time.Second)))
	router.GET("/short", func(c *gin.Context) {
		defer func(writer gin.ResponseWriter) {
			assert.Equal(t, 17, c.Writer.Size())
		}(c.Writer)

		c.JSON(http.StatusOK, gin.H{"hello": "short"})
	})

	req := httptest.NewRequest("GET", "/short", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	result := w.Result()

	assert.Equal(t, http.StatusOK, result.StatusCode)

}
