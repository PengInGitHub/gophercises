package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dir := "./sample"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	counter := 0

	var toRename []string
	for _, file := range files {
		if file.IsDir() {
		} else {
			_, err := match(file.Name(), 4)
			if err == nil {
				counter++
				toRename = append(toRename, file.Name())
			}
		}
	}

	for _, origFileName := range toRename {
		newFileName, err := match(origFileName, counter)
		if err != nil {
			panic(err)
		}
		origPath := filepath.Join(dir, origFileName)
		newPath := filepath.Join(dir, newFileName)
		fmt.Printf("mv %s => %s\n", origPath, newPath)
		err = os.Rename(origPath, newPath)
		if err != nil {
			panic(err)
		}
	}
}

//match returns the new file name
func match(fileName string, total int) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]                      //the last is extension
	tmp := strings.Join(pieces[0:len(pieces)-1], ".") //join the names by seperator
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}
	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil
}
