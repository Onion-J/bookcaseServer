package common

import (
	"io/fs"
	"io/ioutil"
	"log"
)

// FileForEach 遍历指定文件夹下的文件
func FileForEach(fileFullPath string) []fs.FileInfo {
	files, err := ioutil.ReadDir(fileFullPath)
	if err != nil {
		log.Fatal(err)
	}
	var myFile []fs.FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		myFile = append(myFile, file)
	}
	return myFile
}
