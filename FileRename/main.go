package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dir := "sample"

	type file struct {
		path string
		name string
	}

	var toRename []file

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, file{
				path: path,
				name: info.Name(),
			})
		}
		return nil
	})

	for _, f := range toRename {
		fmt.Printf("%q\n", f)
	}

	for _, orig := range toRename {
		var n file
		var err error
		n.name, err = match(orig.name)
		if err != nil {
			fmt.Println("Got an error in mathcing", orig.path, err.Error)
		}
		n.path = filepath.Join(dir, n.name)
		fmt.Printf("mv %s => %s\n", orig.path, n.path)
		err = os.Rename(orig.path, n.path)
		if err != nil {
			fmt.Println("Got an error in renaming", orig.path, err.Error)
		}
	}
}

//match returns the new file name
func match(fileName string) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]                      //the last is extension
	tmp := strings.Join(pieces[0:len(pieces)-1], ".") //join the names by seperator
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our pattern", fileName)
	}
	return fmt.Sprintf("%s - %d.%s", strings.Title(name), number, ext), nil
}
