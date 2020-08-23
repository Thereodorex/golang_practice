package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

var print = fmt.Println

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func inner_recursive(res []string, path string, printFiles bool, level int) ([]string, error) {
	newFiles, err := recursive(path, printFiles)
	if err != nil {
		return nil, err
	}
	for i, value := range newFiles {
		if level == 0 {
			newFiles[i] = "│\t" + value
		} else {
			newFiles[i] = "\t" + value
		}
	}
	res = append(res, newFiles...)
	return res, nil
}

func recursive(path string, printFiles bool) ([]string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, 5)
	if stat.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		filenames := make([]string, 0, len(files))
		for _, file := range files {
			if file.IsDir() {
				filenames = append(filenames, file.Name())
			} else if printFiles {
				size := file.Size()
				var sizeStr string
				if size == 0 {
					sizeStr = "empty"
				} else {
					sizeStr = fmt.Sprintf("%db", file.Size())
				}
				s := fmt.Sprintf("%s (%s)", file.Name(), sizeStr)
				filenames = append(filenames, s)
			}
		}
		sort.Strings(filenames)
		if len(filenames) != 0 {
			for _, file := range filenames[:len(filenames)-1] {
				res = append(res, "├───"+file)
				if file[len(file)-1] != ')' {
					res, err = inner_recursive(res, path+string(os.PathSeparator)+file, printFiles, 0)
				}
				if err != nil {
					return nil, err
				}
			}
			res = append(res, "└───"+filenames[len(filenames)-1])
			file := filenames[len(filenames)-1]
			if file[len(file)-1] != ')' {
				res, err = inner_recursive(res, path+string(os.PathSeparator)+
					filenames[len(filenames)-1], printFiles, 1)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	res, err := recursive(path, printFiles)
	for _, line := range res {
		bytes := []byte(line + "\n")
		out.Write(bytes)
	}
	return err
}

func main() {
	out := os.Stdout
	if len(os.Args) != 2 && (len(os.Args) != 3 || os.Args[2] != "-f") {
		panic("usage: go run tree.go [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
