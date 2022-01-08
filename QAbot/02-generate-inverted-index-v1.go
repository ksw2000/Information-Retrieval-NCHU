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

// - inverted-index/
//		- 0000000-0095999.txt
//		- 0096000-0191999.txt
//		- 0192000-0287999.txt
//		- 0288000-0383999.txt
//		- 0384000-0479999.txt
//		- 0480000-0575999.txt
//		- 0576000-0671999.txt
//		- 0672000-0767999.txt
//		- 0768000-0863999.txt
//		- 0864000-0959999.txt
//		- 0960000-1055999.txt
//		- 1056000-1151999.txt
//		- 1152000-1215638.txt

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

	// (1) read directory and generate fileNames list
	files, err := ioutil.ReadDir("./ws-result")
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}

	// (2) read files by fileNames list and generate a map
	// key: word
	// val: entry including aid and a counter
	invertedIndexList := []map[string]*el{}
	batchSize := 20
	batchCounter := 0

	// record the start aid and end aid
	var startAid, endAid, tmp int

	for i, f := range fileNames {
		fmt.Println("processing", f)

		// (0) if `f` is the start of batch, get startAid
		if i%batchSize == 0 {
			invertedIndexList = append(invertedIndexList, map[string]*el{})
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

		// (3) if `f` is the end of batch, write result to file at other coroutine
		if i%batchSize == batchSize-1 || i == len(fileNames)-1 {
			fmt.Sscanf(f, "output-%d-%d.txt", &tmp, &endAid)
			fmt.Printf("writing ./inverted-index/%07d-%07d.txt at other coroutine. Batch: %d\n", startAid, endAid, batchCounter)
			go func(batchCounterCopy, startAidCopy, endAidCopy int) {
				defer wg.Done()
				defer fmt.Printf("Finish batch: %d\n", batchCounterCopy)

				// open file (create)
				file, err := os.OpenFile(fmt.Sprintf("./inverted-index/%07d-%07d.txt", startAidCopy, endAidCopy), os.O_RDWR|os.O_CREATE, 0755)
				if err != nil {
					log.Print("open file error")
				}
				defer file.Close()

				// fetch invertedIndex that is belong to this coroutine from invertedIndexList
				invertedIndex := invertedIndexList[batchCounterCopy]
				for k, v := range invertedIndex {
					fmt.Fprintf(file, "%s ", k)
					for current := v; current != nil; current = current.next {
						fmt.Fprintf(file, "%d ", current.aid)
					}
					fmt.Fprintf(file, "\n")
					// clear cache
					delete(invertedIndex, k)
				}
			}(batchCounter, startAid, endAid)
			wg.Add(1)
			batchCounter += 1
		}
	}
	// wait all coroutines done
	wg.Wait()
}
