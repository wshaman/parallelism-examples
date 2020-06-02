package read

import (
	"bufio"
	"os"
	"strings"
)

func Do(fname string) ([]Entry, error) {
	f, err := os.Open(fname)
	if err != nil{
		return nil, err
	}
	entries := make([]Entry, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		ln := scanner.Text()
		d := strings.Split(ln, ";")
		entries = append(entries, Entry{
			IP:     d[0],
			Target: d[1],
		})
	}
	return entries, nil
}