package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//Resets back to Original terminal color. Must use at end of program to reset terminal for other applications.
var reset string = "\033[0m"

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

func makeDir(src, dst string) error {
	fmt.Println(cyan, "Scanning For Files and Creating Directory")

	err := filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				os.MkdirAll(strings.ReplaceAll(path, src, dst), os.ModePerm)
			} else {
				destination := finalDestination(path, src, dst)
				if _, err := os.Stat(destination); os.IsNotExist(err) {
					filesList = append(filesList, path)
					//log.Println(reset, "Adding to fileList", path)
				}
			}
			return nil
		})

	if err != nil {
		return err
	}

	return nil
}

func finalDestination(filepath string, sources string, destination string) string {

	//src: fileList | dst: Z:\test | remove_src: C:\Users\magic\Downloads\
	firstPart := strings.ReplaceAll(sources, sources, destination) //path, old, new

	lastPart := strings.ReplaceAll(filepath, sources, "") //Removes the Beginning /home/ of the s path so it can be used to use the correct path for the destination

	//Rewrites it from D:\Skyrim-copy\\ to D:\Skyrim-copy\ removing the extra \ a the end
	newDest := TrimSuffix(firstPart, "\\") + lastPart

	return newDest

}

func copy(id int, jobs <-chan int, results chan<- int, sources string, destination string) {
	//var BUFFERSIZE int64 = 10

	for j := range jobs {

		filepath := filesList[j] //Gets the file that will be moved over by index of the jobs chan

		newDest := finalDestination(filepath, sources, destination)

		stats, err := os.Stat(filepath)

		log.Println(cyan, "Worker ID:", id, green, "Copying:", blue, filepath, purple, "Size", yellow, stats.Size())

		source, err := os.Open(filepath)
		if err != nil {

		}
		defer source.Close()

		/*_, err = os.Stat(newDest)
		if err == nil {
			//log.Println(red, newDest)
		}*/

		destination, err := os.Create(newDest)
		if err != nil {
		}

		defer destination.Close()

		buf := make([]byte, stats.Size())
		for {
			n, err := source.Read(buf)

			if err != nil && err != io.EOF {
			}

			if n == 0 {
				break
			}

			if _, err := destination.Write(buf[:n]); err != nil {
			}
		}

		log.Println(yellow, "Finished Worker:", id)

		results <- j * 2

	}
}

func main() {
	/*if len(os.Args) != 3 {
		fmt.Printf("usage: %s source destination BUFFERSIZE\n", filepath.Base(os.Args[0]))
		return
	}*/

	log.Println(reset)

	//threads := os.Args[1]
	threads, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Invalid buffer size: %q\n", err)
		return
	}

	sources := os.Args[2]
	destination := os.Args[3]

	//Not sure if I need this anymore
	/*res1 := strings.HasSuffix(destination, "\\")
	if !res1 {
		destination = destination + "\\"

		log.Println(red, "Running", destination, reset)
	}*/

	makeDir(sources, destination)

	numJobs := len(filesList)
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for id := 1; id <= int(threads); id++ {
		go copy(id, jobs, results, sources, destination)
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

	fmt.Println(reset, "Done")
}