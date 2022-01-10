package main

import (
	f "fmt"
	"io"
	l "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	//"golang.org/x/crypto/sha3"
)

var r string = "\033[0m" //Resets back to Original terminal color. Must use at end of program to reset terminal for other applications.
var red string = "\033[31m"
var green string = "\033[32m"
var yellow string = "\033[33m"
var blue string = "\033[34m"
var purple string = "\033[35m"
var cyan string = "\033[36m"
var white string = "\033[37m"

var srcFiles []string
var dstFiles []string

var modifyPath string
var sources string
var destination string

func MakeDir() {
	f.Println("Scanning For Files and Creating Directory")

	RenameLater()

	filepath.Walk(sources, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		dstPath := TrimSuffixAndNewDestination(srcPath)
		if _, modPath := os.Stat(dstPath); os.IsNotExist(modPath) { //Checks to see if the destination already exists
			if info.IsDir() { //Checks to see if it's a directory
				os.MkdirAll(dstPath, os.ModePerm) //Makes the directory
			} else { //Runs if it's a file
				srcFiles = append(srcFiles, srcPath)
				dstFiles = append(dstFiles, dstPath)
			}
		}
		return nil
	})
}

func TrimSuffixAndNewDestination(file string) string {
	if hasSuffix := strings.HasSuffix(file, "\\"); hasSuffix {
		file = file[:len(file)-len("\\")]
	}
	return strings.ReplaceAll(file, sources, modifyPath)
}

//TODO Rename and still not sure if I need anymore
func RenameLater() {
	res1 := strings.HasSuffix(destination, "\\")
	if !res1 {
		destination = destination + "\\" //	l.Println(red, "Running", destination, r)
		l.Println("RenameLater")
	}
}

func CheckFileSize(file string) int64 {
	var buff int64
	stats, err := os.Stat(file)
	if err != nil {
		l.Fatalln("func CheckFileSize:", err)
	}
	if stats.Size() > 2000000000 {
		buff = 2000000000
	} else {
		buff = stats.Size()
	}
	return buff
}

func Copy(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		src := srcFiles[j]
		dst := dstFiles[j]

		var buff int64 = CheckFileSize(src)

		l.Println("Worker:", id, "String:", src)

		source, err := os.Open(src)
		if err != nil {
			l.Fatalln("var source := os.Open(file)")
		}
		defer source.Close()

		//srcHash := createHash(file)

		//createHash(file)

		createDestination, err := os.Create(dst)
		if err != nil {
			l.Fatalln("3")
		}
		defer createDestination.Close()

		buf := make([]byte, buff)
		for {
			n, err := source.Read(buf)
			if err != nil && err != io.EOF {
				l.Fatalln("4")
			}
			if n == 0 {
				break
			}
			if _, err := createDestination.Write(buf[:n]); err != nil {
				l.Fatalln("5")
			}
		}
		results <- j * 2
		l.Println("Finished Worker:", r)
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

func MakeWorkers(jobs <-chan int, results chan<- int, threads int) {
	for id := 1; id <= int(threads); id++ {
		go Copy(id, jobs, results)
	}

}

func AddFilesToChannelJobs(jobs chan<- int, numJobs int) {
	for j := 0; j <= numJobs; j++ {
		if j < numJobs { //This if is for making sure that we don't go over the index of the fileList slice
			jobs <- j
		}
	}
	close(jobs)
}

func CheckForResults(results <-chan int, numJobs int) {
	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

func main() {
	if len(os.Args) < 2 {
		f.Printf("Missing \n", filepath.Base(os.Args[0]))
		return
	}

	threads, err := strconv.Atoi(os.Args[1])
	if err != nil {
		f.Printf("Invalid Number of Threads: %q\n", err)
		return
	}
	sources = os.Args[2]
	destination = os.Args[3]
	modifyPath = strings.ReplaceAll(sources, sources, destination) // /home/deathpoolops/test/hello to /new/dst/test/hello

	MakeDir()

	buffSize := len(srcFiles)
	jobs := make(chan int, buffSize)
	results := make(chan int, buffSize)

	if len(srcFiles) != len(dstFiles) {
		panic(len(srcFiles))
	}

	MakeWorkers(jobs, results, threads)

	AddFilesToChannelJobs(jobs, buffSize)

	CheckForResults(results, buffSize)

	f.Println(r, "Done", buffSize)
}
