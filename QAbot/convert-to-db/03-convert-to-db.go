// After word segmentation, we have a lot of files like:

package main

import (
	"bufio"
	"context"
	"database/sql"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type token struct {
	word    string
	wordLen int
	aid     int
}

func main() {
	log.Println("Start")
	// (1) read file and open database
	const file = "../inverted-index/0000000-1215638.txt"
	conn, err := sql.Open("sqlite3", "./inverted-index.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	ctx := context.Background()

	const createTable = `
	BEGIN TRANSACTION;
	CREATE TABLE IF NOT EXISTS "mapping" (
	"word"	TEXT,
	"aids"	TEXT,
	PRIMARY KEY("word")
	);
	COMMIT;
	`

	if _, err := conn.ExecContext(ctx, createTable); err != nil {
		log.Fatal(err)
	}

	// (2) open f
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// (3) read f line by line
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 100*1024*1024)
	scanner.Split(bufio.ScanLines)
	tx, err := conn.BeginTx(ctx, nil)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		aidListNum := []int{}
		// convert to int
		var num int
		for _, aid := range fields[1:] {
			num, _ = strconv.Atoi(aid)
			aidListNum = append(aidListNum, num)
		}
		// sorting
		sort.Ints(aidListNum)
		// convert to string
		aidListString := []string{}
		for _, aid := range aidListNum {
			aidListString = append(aidListString, strconv.Itoa(aid))
		}

		word := fields[0]
		aidList := strings.Join(aidListString, " ")

		if _, err := tx.ExecContext(ctx, "INSERT INTO mapping(word, aids) values(?, ?)", word, aidList); err != nil {
			if err.(sqlite3.Error).Code != sqlite3.ErrConstraint {
				log.Fatal(err)
			}
		}
	}
	tx.Commit()
	log.Println("Done")

	if scanner.Err() != nil {
		log.Printf("Scanner Error: %s\n", scanner.Err())
	}
}
