package worker

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/wshaman/demo-concur/read"
)

type chanFiles chan string
type chanQuit chan bool

var visitors map[string]int64

var entryMutex sync.Mutex

func doEntry (e read.Entry) {
	entryMutex.Lock()
	defer entryMutex.Unlock()
	if _, ok := visitors[e.IP]; !ok {
		visitors[e.IP] = 0
	}
	visitors[e.IP]++
}

func doFile(myID int, inFiles chanFiles, q chanQuit) {
	log.Printf("Worker %d started", myID)
	for {
		select {
			case fPath := <-inFiles :
				log.Printf("Worker %d got %s file", myID, fPath)
				entries, err := read.Do(fPath)
				if err != nil {
					log.Fatal(err)
				}
				for _, entry := range entries {
					doEntry(entry)
				}
			case <- q :
				log.Printf("Worker %d exits", myID)
				return
		}
	}
}

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

func writeFilesToChan(items []string, ch chanFiles) {
	for _, item := range items {
		ch <- item
	}
}

func ReadFilesWorkers () error {
	numWorkers := 8
	visitors = make(map[string]int64)
	fPaths, err := getFileList()
	if err != nil {
		return err
	}
	ch := make(chanFiles)
	q := make(map[int]chanQuit)
	go writeFilesToChan(fPaths, ch)
	for i := 0; i < numWorkers; i++ {
		q[i] = make(chanQuit)
		go doFile(i, ch, q[i])
	}

	time.Sleep(2 * time.Second)
	for i := 0; i < numWorkers; i++ {
		q[i] <- true
	}

	fmt.Printf("Found %d visitors\n", len(visitors))
	return nil
}
