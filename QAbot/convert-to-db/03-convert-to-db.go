// After word segmentation, we have a lot of files like:

// - convert-to-db/
//		- inverted-index.db (About 1.3GB took 6min)

package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	sqlite3 "github.com/mattn/go-sqlite3"
)

const file = "../inverted-index/0000000-1215638.txt"
const db = "./inverted-index.db"

func main() {
	startTime := time.Now()

	// open database
	conn, err := sql.Open("sqlite3", db)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ctx := context.Background()

	// create table
	const createTable = `
	BEGIN TRANSACTION;
	CREATE TABLE IF NOT EXISTS "mapping" (
	"word"	TEXT,
	"num"	INTEGER,
	"aids"	TEXT,
	PRIMARY KEY("word")
	);
	COMMIT;
	`

	if _, err := conn.ExecContext(ctx, createTable); err != nil {
		log.Fatal(err)
	}

	// open inverted table
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// read line by line
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 100*1024*1024)
	scanner.Split(bufio.ScanLines)

	tx, err := conn.BeginTx(ctx, nil)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		aidList := make([]int, 0, len(fields)-2)

		// convert to int
		for _, token := range fields[2:] {
			aid, _ := strconv.Atoi(token)
			aidList = append(aidList, aid)
		}

		// sorting
		sort.Ints(aidList)

		// convert to string
		var aidListString strings.Builder
		for _, aid := range aidList {
			fmt.Fprintf(&aidListString, "%d ", aid)
		}

		word := fields[0]
		num, _ := strconv.Atoi(fields[1])

		if _, err := tx.ExecContext(ctx, "INSERT INTO mapping(word, num, aids) values(?, ?, ?)",
			word, num, aidListString.String()); err != nil {
			if err.(sqlite3.Error).Code != sqlite3.ErrConstraint {
				log.Fatal(err)
			}
		}
	}
	tx.Commit()

	if scanner.Err() != nil {
		log.Printf("Scanner Error: %s\n", scanner.Err())
	}

	elapsed := time.Since(startTime)
	fmt.Println("Done", elapsed)
}
