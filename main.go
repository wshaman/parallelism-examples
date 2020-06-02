package main

import (
	"log"

	"github.com/wshaman/demo-concur/worker"
)

func main() {
	//if err := gorou.ReadFilesGoroutines(); err != nil {
	if err := worker.ReadFilesWorkers(); err != nil {
		log.Fatal(err)
	}
}
