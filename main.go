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


// ballast files
var fre = []*regexp.Regexp{
	regexp.MustCompile(`^\._.*`),
	regexp.MustCompile(`^\.DS_Store$`),
	regexp.MustCompile(`^Thumbs.db$`),
	regexp.MustCompile(`^node_modules$`),
}

// ballast directories
var dre = []*regexp.Regexp{
	//regexp.MustCompile(`^.+/src/.+/node_modules$`),
}

func check(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func isEmptyDir(name string) (bool, error) {
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
	if err != nil {
		fmt.Printf("Skipped: %s: %s\n", path, err)
		return nil
	}
	base := filepath.Base(path)
	if f.IsDir() {
		for _, r := range dre {
			if r.MatchString(base) {
				ballast = append(ballast, path)
			}
		}
	} else {
		for _, r := range fre {
			if r.MatchString(base) {
				ballast = append(ballast, path)
			}
		}
	}
	return nil
}

func findEmptyDirs(path string, f os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("Skipped: %s: %s\n", path, err)
		return nil
	}
	if f.IsDir() {
		empty, err := isEmptyDir(path)
		check(err)
		if empty {
			empties = append(empties, path)
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
	var root string
	flag.Parse()
	if len(flag.Args()) == 1 {
		root = flag.Arg(0)
	} else {
		root, _ = os.Getwd()
	}
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
