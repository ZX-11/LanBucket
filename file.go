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
	Size size
}

var files = make([]*file, 0, 16)

var findFile = make(map[string]*file, 32)

var existFile = make(map[string]int, 32)

func deleteAll() {
	files = files[:0]
	findFile = make(map[string]*file, 32)
	existFile = make(map[string]int, 32)
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
	if times, ok := existFile[s.Name()]; ok {
		existFile[s.Name()]++
		orderedName = strings.TrimSuffix(s.Name(), filepath.Ext(s.Name())) + fmt.Sprintf("(%v)", times) + filepath.Ext(s.Name())
	} else {
		existFile[s.Name()] = 1
	}

	newFile := &file{orderedName, address, size(s.Size())}
	files = append(files, newFile)
	findFile[orderedName] = newFile

	//fmt.Printf("已添加：%v(%v)\n", orderedName, s.Size())
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

type size int64

const (
	_  = 1 << (10 * iota)
	KB // 1024
	MB // 1048576
	GB // 1073741824
	TB // 1099511627776
	PB // 1125899906842624
	EB // 1152921504606846976
)

func (s size) String() string {
	switch {
	case s < KB:
		return fmt.Sprintf("%d B", s)
	case s < MB:
		return fmt.Sprintf("%.2f KB", float64(s)/float64(KB))
	case s < GB:
		return fmt.Sprintf("%.2f MB", float64(s)/float64(MB))
	case s < TB:
		return fmt.Sprintf("%.2f GB", float64(s)/float64(GB))
	case s < PB:
		return fmt.Sprintf("%.2f TB", float64(s)/float64(TB))
	case s < EB:
		return fmt.Sprintf("%.2f PB", float64(s)/float64(PB))
	}
	return fmt.Sprintf("%.2f EB", float64(s)/float64(EB))
}
