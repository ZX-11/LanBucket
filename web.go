package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var localAddr = func() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return []byte{127, 0, 0, 1}
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}()

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
	r.GET("/file", func(c *gin.Context) {
		name := c.Query("name")
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
