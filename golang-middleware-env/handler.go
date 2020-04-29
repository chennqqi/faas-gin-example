package function

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	engine = gin.Default()
	server = &Server{}
)

func init() {
	r := engine
	r.GET("/dump", server.dump)
	r.GET("/", server.dump)
	r.GET("/dump/dump", server.dump)
	r.GET("/id/:id/name/:name", server.dumpParam)
	r.GET("/vid/:id/vname/*name", server.dumpParam)
	r.GET("/vvid/:id/*name", server.dumpParam)
	r.PUT("/log", server.log)
}

type Dump struct {
	Header http.Header `json:"headers"`
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Query  string      `json:"query"`
	Host   string      `json:"host"`
}

type Server struct {
}

func (self *Server) log(c *gin.Context) {
	if c.Request.Body != nil {
		io.Copy(c.Writer, c.Request.Body)
		defer c.Request.Body.Close()
	} else {
		c.String(200, "no body")
	}
}

func (self *Server) dump(c *gin.Context) {
	var d Dump
	r := c.Request
	d.Header = r.Header
	d.URL = r.URL.String()
	d.Method = r.Method
	d.Query = r.URL.RawQuery
	d.Host = r.Host
	c.JSON(200, d)
}

func (self *Server) dumpParam(c *gin.Context) {
	name := c.Param("name")
	id := c.Param("id")
	c.JSON(200, map[string]interface{}{
		"name": name,
		"id":   id,
	})
}

func Handle(w http.ResponseWriter, r *http.Request) {
	engine.ServeHTTP(w, r)
}
