// - word-frequency-table/
//		- word-frequency.db (About 340MB took 5min)

package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
)

const dir = "../ws-result"

func main() {
	startTime := time.Now()

	// STEP1: read directory and generate fileNames list
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	fileNames := []string{}

	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	sort.Strings(fileNames)

	// STEP2: connect to database
	ctx := context.Background()
	conn, err := sql.Open("sqlite3", "./word-frequency.db")
	const createTable = `
	BEGIN TRANSACTION;
	CREATE TABLE IF NOT EXISTS "frequency" (
	"word"	TEXT,
	"num"	INTEGER,
	PRIMARY KEY("word")
	);
	COMMIT;
	`

	if _, err := conn.ExecContext(ctx, createTable); err != nil {
		log.Fatal(err)
	}

	// STEP3 read files from fileNames list and generate a map
	// key: word val: number of article
	frequency := map[string]int{}

	for _, f := range fileNames {
		fmt.Println("processing", f)

		// open file
		file, err := os.Open(path.Join(dir, f))
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
			for i, word := range strings.Fields(scanner.Text()) {
				if i != 0 {
					frequency[word] += 1
				} else {
					fmt.Sscanf(word, "%d", &aid)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// STEP4: write to database.
	tx, err := conn.BeginTx(ctx, nil)
	for word, num := range frequency {
		if _, err := tx.ExecContext(ctx, "INSERT INTO frequency(word, num) values(?, ?)", word, num); err != nil {
			if err.(sqlite3.Error).Code != sqlite3.ErrConstraint {
				log.Fatal(err)
			}
		}
	}
	tx.Commit()

	elapsed := time.Since(startTime)
	fmt.Println("Done", elapsed)
}
