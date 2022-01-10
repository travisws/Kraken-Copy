package main_test

import (
	"testing"
	"io"
	"os"
	"path/filepath"
	"strings"
	l "log"

)

var filesList []string

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func MakeDir(src, dst string) {
	filepath.Walk(src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				l.Println(path)
				os.MkdirAll(strings.ReplaceAll(path, src, dst), os.ModePerm)
			} else {
				destination := FinalDestination(path, src, dst)
				if _, err := os.Stat(destination); os.IsNotExist(err) {
					filesList = append(filesList, path)
					l.Println("Adding to fileList", path)
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

func Copy(id int, jobs <-chan int, results chan<- int, sources string, destination string) {
	for j := range jobs {
		var buff int64 = 10

		filepath := filesList[j] //Gets the file that will be moved over by index of the jobs chan

		newDest := FinalDestination(filepath, sources, destination)

		stats, err := os.Stat(filepath)

		if stats.Size() > 2000000000 {
			buff = 2000000000
		} else {
			buff = stats.Size()
		}


		source, err := os.Open(filepath)
		if err != nil {

		}
		defer source.Close()

		//srcHash := createHash(filepath)

		//createHash(filepath)

		destination, err := os.Create(newDest)
		if err != nil {
		}
		defer destination.Close()

		buf := make([]byte, buff)
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

		results <- j * 2

	}
}

func TestMain(t *testing.T){
	sources := "D:\\Minecraft"
	destination := "E:\\Minecraft"
	threads := 5

	MakeDir(sources, destination)

	numJobs := len(filesList)
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)


	for id := 10; id <= int(threads); id++ {
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

}