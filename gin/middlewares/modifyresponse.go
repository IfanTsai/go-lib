package middlewares

import (
	"bytes"
	"strconv"

	"github.com/IfanTsai/go-lib/set"

	"github.com/gin-gonic/gin"
)

// ModifyResponse modify response body in cb before returning to client,
// can filter url through urls parameter, does not filter if urls is nil
func ModifyResponse(urls *set.Set, cb func(c *gin.Context, url string, body *bytes.Buffer)) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if urls != nil && !urls.Contains(url) {
			return
		}

		w := NewResponseBufferWriter(c.Writer)
		c.Writer = w

		c.Next()

		cb(c, url, w.body)
		w.Flush()
	}
}

type ResponseBufferWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	status  int
	flushed bool
}

func NewResponseBufferWriter(w gin.ResponseWriter) *ResponseBufferWriter {
	return &ResponseBufferWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (w *ResponseBufferWriter) Write(buf []byte) (int, error) {
	return w.body.Write(buf)
}

func (w *ResponseBufferWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *ResponseBufferWriter) Written() bool {
	return w.body.Len() > 0
}

func (w *ResponseBufferWriter) WriteHeader(status int) {
	w.status = status
}

func (w *ResponseBufferWriter) Status() int {
	return w.status
}

func (w *ResponseBufferWriter) Size() int {
	return w.body.Len()
}

func (w *ResponseBufferWriter) Flush() {
	if w.flushed {
		return
	}

	w.ResponseWriter.WriteHeader(w.status)
	if w.body.Len() > 0 {
		w.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(w.Size()))
		if _, err := w.ResponseWriter.Write(w.body.Bytes()); err != nil {
			panic(err)
		}

		w.body.Reset()
	}

	w.flushed = true
}
