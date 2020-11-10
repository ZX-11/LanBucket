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

var stdin = bufio.NewReader(os.Stdin)

var sigs = make(chan os.Signal, 1)

func main() {
	fmt.Println("把需要添加的文件或目录拖拽进来即可添加")
	fmt.Printf("访问地址：http://%v:18800\n", getIP())
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
				"files": files,
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
		r.Run(":18800")
	}()

	signal.Notify(sigs, os.Interrupt, os.Kill)

	<-sigs
}

func readLine() string {
	buf, _, _ := stdin.ReadLine()
	return string(buf)
}

func getIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP
}
