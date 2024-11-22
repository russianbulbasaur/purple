package rdb

import (
	"log"
	"os"
	"path/filepath"
)

type RDBFile struct {
	file       *os.File
	fullPath   string
	dir        string
	dbFileName string
	size       int64
}

func NewRDBFile(filename string, dir string) *RDBFile {
	_, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}
	filePath := filepath.Join(dir, filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	info, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return &RDBFile{
		file,
		filePath,
		dir,
		filename,
		info.Size(),
	}
}

func (file *RDBFile) GetDir() string {
	return file.dir
}

func (file *RDBFile) GetDBFileName() string {
	return file.dbFileName
}
