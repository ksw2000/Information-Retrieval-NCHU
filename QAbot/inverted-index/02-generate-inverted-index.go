// After word segmentation, we have a lot of files which are consisted of:

// aid1 word1 word2 word3
// aid2 word1 word2 word3 word4
// aid3 word1

// Try to create inverted-index files
// version 3 writes frequency into files

// CPU: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
// RAM: 19.9 GB (Need RAM about 7GB)
// Disk: CT500MX500SSD

// - inverted-index/
//		- 0000000-1215638.txt (About 1.1GB took 10min)

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const readDir = "../ws-result-level1/"
const outputDir = "./"

func main() {
	startTime := time.Now()

	// STEP1: read directory and generate fileNames list
	files, err := ioutil.ReadDir(readDir)
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	sort.Strings(fileNames)

	// STEP2: read files from fileNames list and generate a map
	// key: word val: entry including aid
	invertedIndex := map[string][]int{}
	for _, f := range fileNames {
		fmt.Println("processing", f)

		// open f
		file, err := os.Open(path.Join(readDir, f))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// read f line by line
		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			var aid int

			// read every words in single line
			dict := map[string]interface{}{}
			for i, word := range strings.Fields(scanner.Text()) {
				if i != 0 {
					dict[word] = nil
				} else {
					fmt.Sscanf(word, "%d", &aid)
				}
			}

			// generate invertedIndex
			for word := range dict {
				invertedIndex[word] = append(invertedIndex[word], aid)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// STEP3: Write result
	// generate file name
	var tmp, endAid int
	fmt.Sscanf(fileNames[len(fileNames)-1], "output-%d-%d.txt", &tmp, &endAid)
	generateFile := path.Join(outputDir, fmt.Sprintf("%07d-%07d.txt", 0, endAid))
	fmt.Printf("writing %s\n", generateFile)

	// open file
	file, err := os.OpenFile(generateFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Print("open file error")
	}
	defer file.Close()

	for k, v := range invertedIndex {
		fmt.Fprintf(file, "%s %d ", k, len(v))
		for _, aid := range v {
			fmt.Fprintf(file, "%d ", aid)
		}
		fmt.Fprintf(file, "\n")
	}

	elapsed := time.Since(startTime)
	fmt.Println("Done", elapsed)
}
