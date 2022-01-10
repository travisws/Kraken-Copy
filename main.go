package main

import (
	//"bytes"
	f "fmt"
	"io"
	l "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	//"golang.org/x/crypto/sha3"
)

//Resets back to Original terminal color. Must use at end of program to reset terminal for other applications.
var r string = "\033[0m"

var red string = "\033[31m"
var green string = "\033[32m"
var yellow string = "\033[33m"
var blue string = "\033[34m"
var purple string = "\033[35m"
var cyan string = "\033[36m"
var white string = "\033[37m"

var filesList []string

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func MakeDir(src, dst string) {
	f.Println(cyan, "Scanning For Files and Creating Directory")

	filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				//l.Println(green, path)
				os.MkdirAll(strings.ReplaceAll(path, src, dst), os.ModePerm)
			} else {
				destination := FinalDestination(path, src, dst)
				if _, err := os.Stat(destination); os.IsNotExist(err) {
					filesList = append(filesList, path)
					//l.Println(r, "Adding to fileList", path)
				}
			}
			return nil
		})
}

func FinalDestination(filepath string, sources string, destination string) string {
	//src: fileList | dst: Z:\test | remove_src: C:\Users\magic\Downloads\
	firstPart := strings.ReplaceAll(sources, sources, destination) //path, old, new

	lastPart := strings.ReplaceAll(filepath, sources, "") //Removes the Beginning /home/ of the s path so it can be used to use the correct path for the destination

	//Rewrites it from D:\Skyrim-copy\\ to D:\Skyrim-copy\ removing the extra \ a the end
	newDest := TrimSuffix(firstPart, "\\") + lastPart

	return newDest

}

var tJobs int64 = 0

func Copy(id int, jobs <-chan int, results chan<- int, sources string, destination string) {
	for j := range jobs {
		var buff int64 = 10

		filepath := filesList[j] //Gets the file that will be moved over by index of the jobs chan

		newDest := FinalDestination(filepath, sources, destination)

		l.Println(r, "filepath:", filepath, "\n", "newdest", newDest, r)

		stats, err := os.Stat(filepath)
		if err != nil {
			l.Fatalln("1")
		}

		l.Println(cyan, "Worker ID:", id, "Copying:", filepath, "Size", stats.Size(), r)

		if stats.Size() > 2000000000 {
			buff = 2000000000
		} else {
			buff = stats.Size()
		}

		//l.Println(red, "buff", buff, r)

		source, err := os.Open(filepath)
		if err != nil {
			l.Fatalln("2")
		}
		defer source.Close()

		//srcHash := createHash(filepath)

		//createHash(filepath)

		destination, err := os.Create(newDest)
		if err != nil {
			l.Fatalln("3")
		}
		defer destination.Close()

		buf := make([]byte, buff)
		for {
			n, err := source.Read(buf)

			if err != nil && err != io.EOF {
				l.Fatalln("4")
			}

			if n == 0 {
				break
			}

			if _, err := destination.Write(buf[:n]); err != nil {
				l.Fatalln("5")
			}
		}

		tJobs++

		results <- j * 2

		l.Println(yellow, "Finished Worker:", id, r)

	}
}

/*
func createHash(src string)[]byte {
	input := strings.NewReader(src)
	hash := sha3.New512()

	if _, err := io.Copy(hash, input); err != nil {
		l.Println(red, "CAN NOT READ FILE:", src, "\n ERROR:", err)
	}

	sum := hash.Sum(nil)

	l.Println(red, "HASH:", sum)

	return sum
}*/

/*func verifyHash() {

}*/

func main() {
	/*if len(os.Args) != 3 {
		f.Printf("usage: %s source destination BUFFERSIZE\n", filepath.Base(os.Args[0]))
		return
	}*/

	//l.Println(r)

	threads, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		f.Printf("Invalid buffer size: %q\n", err)
		return
	}
	sources := os.Args[2]
	destination := os.Args[3]

	//Not sure if I need this anymore
	res1 := strings.HasSuffix(destination, "\\")
	if !res1 {
		destination = destination + "\\"

		//	l.Println(red, "Running", destination, r)
	}

	/*path, err := os.Open(sources)
	if err != nil {
		// handle the error and return
		l.Fatalln("Hello")
	}

	pathInfo, err := path.Stat()
	if err != nil {
		// error handling
		l.Fatalln("Hello")
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if pathInfo.IsDir() {
		// file is a directory
		l.Fatalln("Hello")

	} else {
		// file is not a directory
	}
	defer path.Close()*/

	MakeDir(sources, destination)

	numJobs := len(filesList)
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for id := 1; id <= int(threads); id++ {
		go Copy(id, jobs, results, sources, destination)
	}

	for j := 0; j <= numJobs; j++ {
		if j < numJobs { //This if is for making sure that we don't go over the index of the fileList slice
			jobs <- j
		}
	}

	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}

	f.Println(r, "Done", tJobs)
}