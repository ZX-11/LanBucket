package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var tmpl string

var ipv4Addr = func() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}()

var ipv6Addr = func() string {
	conn, err := net.Dial("udp", "[2001:4860:4860::8888]:80")
	if err != nil {
		return "[::1]"
	}
	defer conn.Close()
	return "[" + conn.LocalAddr().(*net.UDPAddr).IP.String() + "]"
}()

func webAddr() string {
	if settings["EnableIPv6"].(bool) {
		return fmt.Sprintf("http://%v%v", ipv6Addr, port)
	} else {
		return fmt.Sprintf("http://%v%v", ipv4Addr, port)
	}
}

func webService() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	t, _ := template.New("tmpl").Parse(tmpl)
	r.SetHTMLTemplate(t)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "tmpl", gin.H{
			"files":        files,
			"enableUpload": settings["EnableUpload"],
		})
	})
	r.GET("/file/:name", func(c *gin.Context) {
		name := c.Param("name")
		if file, ok := findFile[name]; ok {
			c.File(file.Addr)
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})
		}
	})
	r.GET("/dl/:name", func(c *gin.Context) {
		name := c.Param("name")
		if file, ok := findFile[name]; ok {
			c.FileAttachment(file.Addr, file.Name)
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})
		}
	})
	r.POST("/upload", func(c *gin.Context) {
		if settings["EnableUpload"].(bool) {
			f, err := c.FormFile("选择文件")
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": err,
				})
			} else {
				os.Mkdir("upload", 0666)
				c.SaveUploadedFile(f, `./upload/`+f.Filename)
				if err := add(`./upload/` + f.Filename); err != nil {
					log.Println(err)
				}
				fileUpload <- struct{}{}
				c.Redirect(http.StatusMovedPermanently, "/")
			}
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})
		}
	})
	r.Run(port)
}
