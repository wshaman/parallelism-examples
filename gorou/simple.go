package gorou

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"sync/atomic"

	"github.com/wshaman/demo-concur/read"
)

var counter int64

func getFileList() ([]string, error) {
	p := path.Join(os.TempDir(), "logsFolder")
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}
	fPaths := make([]string, 0, len(files))
	for _, v := range files {
		if v.IsDir() {
			continue
		}
		fPaths = append(fPaths, path.Join(p, v.Name()))
	}
	return fPaths, nil
}

func ReadFiles () error {
	cnt := 0
	fPaths, err := getFileList()
	if err != nil {
		return err
	}
	for _, v := range fPaths {
		entries, err := read.Do(v)
		if err != nil {
			return err
		}
		cnt += len(entries)
	}
	fmt.Printf("%d records found\n", cnt)
	return nil
}

func doFile(fPath string, wg *sync.WaitGroup) {
	defer wg.Done()
	entries, err := read.Do(fPath)
	if err != nil {
		log.Fatal(err)
	}
	atomic.AddInt64(&counter, int64(len(entries)))
}

func ReadFilesGoroutines () error {
	fPaths, err := getFileList()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(fPaths))
	for _, v := range fPaths {
		go doFile(v, &wg)
	}
	wg.Wait()
	log.Printf("%d records found\n", counter)
	return nil
}

