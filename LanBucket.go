package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

var stdin = bufio.NewScanner(os.Stdin)

var sigs = make(chan os.Signal, 1)

func main() {
	fmt.Println("把需要添加的文件或目录拖拽进来即可添加")
	fmt.Printf("访问地址：http://%v:18800\n", getIP())
	fmt.Println(`您还可以使用“enable/disable upload”命令打开或关闭上传功能`)
	go loadFiles()

	f, _ := os.Create("gin.log")
	defer f.Close()

	go func() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.MultiWriter(f)
		r := gin.Default()
		t, _ := template.New("tmpl").Parse(tmpl)
		r.SetHTMLTemplate(t)
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "tmpl", gin.H{
				"files":        files,
				"enableUpload": enableUpload,
			})
		})
		r.GET("/file", func(c *gin.Context) {
			name := c.Query("name")
			if file, ok := getFileByName[name]; ok {
				c.FileAttachment(file.Addr, file.Name)
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Not Found",
				})
			}
		})
		r.POST("/upload", func(c *gin.Context) {
			if enableUpload {
				f, err := c.FormFile("选择文件")
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{
						"error": err,
					})
				} else {
					c.SaveUploadedFile(f, `./upload/`+f.Filename)
					fmt.Println("接收到上传文件：" + f.Filename)
					if err := add(`./upload/` + f.Filename); err != nil {
						fmt.Println(err)
					}
					c.Redirect(http.StatusMovedPermanently, "/")
				}
			} else {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Not Found",
				})
			}
		})
		r.Run(":18800")
	}()

	signal.Notify(sigs, os.Interrupt, os.Kill)

	<-sigs
}

func getIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}
