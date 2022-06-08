package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

func CsvReader(f *os.File, extractCh chan []string) {
	fmt.Println("csvReader function started")
	r := csv.NewReader(f)
	
	// Go trough every record we have.
	for record, err := r.Read(); err == nil; record, err = r.Read() {
		// Send it to the channel.
		extractCh <- record
	}
}

func JsonReader(f *os.File, extractCh chan []string) {
	//TODO
}

func TxtReader(f *os.File, extractCh chan []string) {
	//TODO
}

func FlatReader(f *os.File, extractCh chan []string) {
	//TODO
}