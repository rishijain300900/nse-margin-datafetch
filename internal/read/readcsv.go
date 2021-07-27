package readcsv

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadCsv(name string) [][]string {
	f, err := os.Open(filepath.Join("Data", filepath.Base(name)))
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", name, err.Error())
	}
	defer f.Close()
	row1, err := bufio.NewReader(f).ReadSlice('\n')
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Seek(int64(len(row1)), io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}
