package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var stdin = bufio.NewScanner(os.Stdin)

func cmdService() {
	fmt.Println("把需要添加的文件或目录拖拽进来即可添加")
	fmt.Printf("访问地址：http://%v%v\n", localAddr, port)
	fmt.Println(`您还可以使用“enable/disable upload”命令打开或关闭上传功能`)
	for stdin.Scan() {
		line := strings.Trim(stdin.Text(), `"`)
		switch {
		case line == "" || line == `\` || line == `/`:
			// do nothing
		case line == "enable upload":
			settings["EnableUpload"] = true
			fmt.Println("已开启上传功能")
		case line == "disable upload":
			settings["EnableUpload"] = false
			fmt.Println("已关闭上传功能")
		default:
			if err := add(line); err != nil {
				fmt.Println(err)
			}
		}
	}
}
