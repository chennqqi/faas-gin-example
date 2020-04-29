package faasginwrapper

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openfaas-incubator/go-function-sdk"
)

var (
	DefaultRouteKey  = "X-Golang-Http-Path"
	DefaultRouteHost = "golangginfakehost"
)

type ResponseWriter struct {
	err error
	buf bytes.Buffer
	rw  handler.Response
}

func (self *ResponseWriter) Header() http.Header {
	return self.rw.Header
}

func (self *ResponseWriter) Write(dat []byte) (int, error) {
	return self.buf.Write(dat)
}

func (self *ResponseWriter) WriteHeader(statusCode int) {
	self.rw.StatusCode = statusCode
}

type ginHandle struct {
	engine *gin.Engine
}

func (self *ginHandle) ServeHandle(faasReq handler.Request) (handler.Response, error) {
	engine := self.engine

	var body *bytes.Buffer
	if len(faasReq.Body) > 0 {
		body = bytes.NewBuffer(faasReq.Body)
	}
	route := faasReq.Header.Get(DefaultRouteKey)
	url := fmt.Sprintf("http://%s%s?%s", DefaultRouteHost, route, faasReq.QueryString)

	req, err := http.NewRequest(faasReq.Method, url, body)
	if err != nil {
		return handler.Response{}, err
	}

	req.Header = faasReq.Header
	var rw ResponseWriter

	engine.ServeHTTP(&rw, req)
	rw.rw.Body = rw.buf.Bytes()
	return rw.rw, nil
}

//golang-http gin Wrapper
func GinFaasHandler(engine *gin.Engine) func(faasReq handler.Request) (handler.Response, error) {
	var h = &ginHandle{
		engine: engine,
	}
	return h.ServeHandle
}

//golang-middleware gin Wrapper
func GinFaasMiddleware(engine *gin.Engine) func(w http.ResponseWriter, r *http.Request) {
	return engine.ServeHTTP
}
