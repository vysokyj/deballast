package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var ballast = make([]string, 0)
var empties = make([]string, 0)

var regexps = []*regexp.Regexp{
	regexp.MustCompile(`^\._.*`),
	regexp.MustCompile(`^\.DS_Store$`),
	regexp.MustCompile(`^Thumbs.db$`),
}

func check(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func IsEmptyDir(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func findBallast(path string, f os.FileInfo, err error) error {
	//fmt.Printf("Visited: %s\n", path)
	fileInfo, err := os.Stat(path)
	check(err)
	if !fileInfo.IsDir() {
		base := filepath.Base(path)
		for _, regexp := range regexps {
			if regexp.MatchString(base) {
				ballast = append(ballast, path)
				//fmt.Printf("Ballast: %s\n", path)
			}
		}
	}
	return nil
}

func findEmptyDirs(path string, f os.FileInfo, err error) error {
	//fmt.Printf("Visited: %s\n", path)
	fileInfo, err := os.Stat(path)
	check(err)
	if fileInfo.IsDir() {
		empty, err := IsEmptyDir(path)
		check(err)
		if empty {
			empties = append(empties, path)
			//fmt.Printf("Empty: %s\n", path)
		}
	}
	return nil
}

func deleteAll(paths []string) {
	for _, path := range paths {
		err := os.Remove(path)
		check(err)
		fmt.Printf("Removed: %s\n", path)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Missing directory argument!")
		os.Exit(1)
	}
	root := flag.Arg(0)
	err := filepath.Walk(root, findBallast)
	check(err)
	deleteAll(ballast)
	ballast = nil
	err = filepath.Walk(root, findEmptyDirs)
	check(err)
	for len(empties) > 0 {
		deleteAll(empties)
		empties = nil
		if exists(root) {
			err = filepath.Walk(root, findEmptyDirs)
			check(err)
		}
	}
}
