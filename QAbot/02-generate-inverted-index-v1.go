// After word segmentation, we have a lot of files like:

// aid1 word1 word2 word3
// aid2 word1 word2 word3 word4
// aid3 word1

// Try to create inverted-index files
// version 1 does not write frequency into files

// About 6 mins
// CPU: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz
// RAM: 19.9 GB
// Disk: CT500MX500SSD

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

type el struct {
	aid  int
	num  int
	next *el
}

func main() {
	// wating threads
	var wg sync.WaitGroup

	// (1) read directory
	files, err := ioutil.ReadDir("./ws-result")
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	// (2) read files
	// generate map
	invertedIndexList := []map[string]*el{}
	batchSize := 20
	batchCounter := 0
	// record the start aid and end aid
	var startAid, endAid, tmp int
	skip := 0
	for i, f := range fileNames {
		if i < skip {
			continue
		}
		fmt.Println("processing", f)
		if i%batchSize == 0 {
			invertedIndexList = append(invertedIndexList, map[string]*el{})
			fmt.Sscanf(f, "output-%d-%d.txt", &startAid, &tmp)
		}
		file, err := os.Open("./ws-result/" + f)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// read one file line by line
		scanner := bufio.NewScanner(file)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		// read one line
		for scanner.Scan() {
			var aid int
			// read every word
			frequency := map[string]int{}
			for i, word := range strings.Fields(scanner.Text()) {
				if i != 0 {
					frequency[word] += 1
				} else {
					fmt.Sscanf(word, "%d", &aid)
				}
			}

			for word, num := range frequency {
				e, ok := invertedIndexList[batchCounter][word]
				if ok {
					e.next = &el{
						aid:  aid,
						num:  num,
						next: e.next,
					}
				} else {
					invertedIndexList[batchCounter][word] = &el{
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

		if i%batchSize == batchSize-1 || i == len(fileNames)-1 {
			fmt.Sscanf(f, "output-%d-%d.txt", &tmp, &endAid)
			fmt.Printf("writing ./inverted-index/%07d-%07d.txt at other thread. Batch: %d\n", startAid, endAid, batchCounter)
			// write file in another thread
			go func(batchCounterCopy, startAidCopy, endAidCopy int) {
				defer wg.Done()
				defer fmt.Printf("Finish batch: %d\n", batchCounterCopy)
				file, err := os.OpenFile(fmt.Sprintf("./inverted-index/%07d-%07d.txt", startAidCopy, endAidCopy), os.O_RDWR|os.O_CREATE, 0755)
				if err != nil {
					log.Print("open file error")
				}
				defer file.Close()
				mapForMe := invertedIndexList[batchCounterCopy]
				for k, v := range mapForMe {
					fmt.Fprintf(file, "%s ", k)
					for current := v; current != nil; current = current.next {
						fmt.Fprintf(file, "%d ", current.aid)
					}
					fmt.Fprintf(file, "\n")
					// clear cache
					delete(mapForMe, k)
				}
			}(batchCounter, startAid, endAid)
			wg.Add(1)
			batchCounter += 1
		}
	}
	// wait threads
	wg.Wait()
}

//
// - inverted-index/
// 		- 0000000-0095999.txt
// 		- 0096000-0191999.txt
// 		- ...
