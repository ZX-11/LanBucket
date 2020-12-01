package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type file struct {
	Name string
	Addr string
	Size int64
}

var files = make([]*file, 0, 16)

var getFileByName = make(map[string]*file, 32)

var exist = make(map[string]int, 32)

var enableUpload = false

func loadFiles() {
	for stdin.Scan() {
		line := strings.Trim(stdin.Text(), `"`)
		switch {
		case line == "" || line == `\` || line == `/`:
			// do nothing
		case line == "enable upload":
			enableUpload = true
			os.Mkdir("upload", 0666)
			fmt.Println("已开启上传功能")
		case line == "disable upload":
			enableUpload = false
			fmt.Println("已关闭上传功能")
		default:
			if err := add(line); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func add(address string) error {
	s, err := os.Stat(address)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return addDir(address)
	}

	orderedName := s.Name()
	if times, ok := exist[s.Name()]; ok {
		exist[s.Name()]++
		orderedName = strings.TrimSuffix(s.Name(), filepath.Ext(s.Name())) + fmt.Sprintf("(%v)", times) + filepath.Ext(s.Name())
	} else {
		exist[s.Name()] = 1
	}

	newFile := &file{orderedName, address, s.Size()}

	files = append(files, newFile)
	getFileByName[orderedName] = newFile

	fmt.Printf("已添加：%v(%vB)\n", orderedName, s.Size())
	return nil
}

func addDir(dir string) error {
	files, err := filepath.Glob(dir + "/*")
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := add(f); err != nil {
			return err
		}
	}
	return nil
}
