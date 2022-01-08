// After word segmentation, we have a lot of files like:

// aid1 word1 word2 word3
// aid2 word1 word2 word3 word4
// aid3 word1

// Try to create inverted-index files
// version 2 does not write frequency into files, as well as version1
// version 2 neither uses multiple goroutines nor uses batch

// CPU: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
// RAM: 19.9 GB (記憶體這麼大94爽 ^_^)
// Disk: CT500MX500SSD

// - inverted-index/
//		- 0000000-1215638.txt (About 1GB)

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type el struct {
	aid  int
	num  int
	next *el
}

func main() {
	// (1) read directory and generate fileNames list
	files, err := ioutil.ReadDir("./ws-result")
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	sort.Strings(fileNames)

	// (2) read files from fileNames list and generate a map
	// key: word
	// val: entry including aid and a counter
	invertedIndex := map[string]*el{}
	// record the start aid and end aid
	var startAid, endAid, tmp int

	for i, f := range fileNames {
		fmt.Println("processing", f)

		if i == 0 {
			fmt.Sscanf(f, "output-%d-%d.txt", &startAid, &tmp)
		}

		// (1) open f
		file, err := os.Open("./ws-result/" + f)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// (2) read f line by line
		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			var aid int

			// read every words in single line
			frequency := map[string]int{}
			for i, word := range strings.Fields(scanner.Text()) {
				if i != 0 {
					frequency[word] += 1
				} else {
					fmt.Sscanf(word, "%d", &aid)
				}
			}

			// generate invertedIndex
			for word, num := range frequency {
				e, ok := invertedIndex[word]
				if ok {
					e.next = &el{
						aid:  aid,
						num:  num,
						next: e.next,
					}
				} else {
					invertedIndex[word] = &el{
						aid:  aid,
						num:  num,
						next: nil,
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// (3) write result to file at other coroutine
		if i == len(fileNames)-1 {
			fmt.Sscanf(f, "output-%d-%d.txt", &tmp, &endAid)
			fmt.Printf("writing ./inverted-index/%07d-%07d.txt at other coroutine.\n", startAid, endAid)
			// open file (create) About 1GB
			file, err := os.OpenFile(fmt.Sprintf("F:/%07d-%07d.txt", startAid, endAid), os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				log.Print("open file error")
			}
			defer file.Close()

			for k, v := range invertedIndex {
				fmt.Fprintf(file, "%s ", k)
				for current := v; current != nil; current = current.next {
					fmt.Fprintf(file, "%d ", current.aid)
				}
				fmt.Fprintf(file, "\n")
				// clear cache
				delete(invertedIndex, k)
			}
		}
	}
}
