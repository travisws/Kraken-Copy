package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)
//Resets back to Original terminal color. Must use at end of program to reset terminal for other applications.
var colorReset string = "\033[0m"

var colorRed string = "\033[31m"
var colorGreen string = "\033[32m"
var colorYellow string = "\033[33m"
var colorBlue string = "\033[34m"
var colorPurple string = "\033[35m"
var colorCyan string = "\033[36m"
var colorWhite string = "\033[37m"

var filesList []string

func makeDir(src, dst string) error {
	var i int = 0

	err := filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				dir := os.MkdirAll(strings.ReplaceAll(path, src, dst), os.ModePerm)
				if err == nil {
					return err
				}
				fmt.Println(colorCyan,"Creating Directory", colorGreen, dir)
			} else {
				filesList = append(filesList, path)
				i += 1
			}
			return nil
		})

	if err != nil {
		return err
	}

	return nil
}

func copy(src, dst string, BUFFERSIZE int64) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}	
	return err
}

func TrimSuffix(s, suffix string) string {
    if strings.HasSuffix(s, suffix) {
        s = s[:len(s)-len(suffix)]
    }
    return s
}


var BUFFERSIZE int64

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("usage: %s source destination BUFFERSIZE\n", filepath.Base(os.Args[0]))
		return
	}

	sources := os.Args[1]
	destination := os.Args[2]
	BUFFERSIZE, err := strconv.ParseInt(os.Args[3], 10, 64)	
	if err != nil {
		fmt.Printf("Invalid buffer size: %q\n", err)
		return
	}

	fmt.Println(colorCyan, "Scanning For Files and Creating Directory")

	res1 := strings.HasSuffix(destination, "\\")
	if !res1 {
		destination = destination + "\\"
	}

	//Runs the func makeDir with 2 strings input
	makeDir(sources, destination)

	totalNumOfFiles := len(filesList)

	test := len(filesList) / 2

	for _, source := range filesList {
		lastPart := strings.ReplaceAll(source, sources, "") 
		//src: fileList | dst: Z:\test | remove_src: C:\Users\magic\Downloads\
		firstPart := strings.ReplaceAll(sources, sources, destination) //path, old, new

		//Rewrites it from D:\Skyrim-copy\\ to D:\Skyrim-copy\ removing the extra \ a the end
		newDest := TrimSuffix(firstPart, "\\") + lastPart

		err = copy(source, newDest, BUFFERSIZE)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		} else {
			fmt.Println(colorCyan,"File Number",  "of", totalNumOfFiles, colorGreen, "Coping", colorBlue, source, colorPurple, "To", colorYellow, newDest)
		}
		fmt.Println(test)
	}
	fmt.Println(colorReset, "Done")
		//Come back to 
}