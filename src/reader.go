package src

import (
	"encoding/json"
	"os"
	"proj2/png"
	"sync"
)

// struct for decoding image process request
type Request struct{
	InPath string
	OutPath string
	Effects []string
}

// struct for a pre-processed request
type FilterRequest struct{
	img *png.Image
	req Request
}

// A reader reads from os.stdin and create its own work pipeline for image process
func readAndProcess(wg *sync.WaitGroup, id int, cond *sync.Cond, dec *json.Decoder, n int){
	defer wg.Done()
	done := make(chan bool, 1)
	filterReq := make(chan FilterRequest, 10)

	var tmp Request
	for{
		cond.L.Lock()
		err := dec.Decode(&tmp)
		cond.L.Unlock()
		if err != nil {
			close(filterReq)
			return
		}

		filePath := tmp.InPath
		pngImg, err := png.Load(filePath)
		if err != nil {
			panic(err)
		}

		var tmpReq FilterRequest
		tmpReq.img = pngImg
		tmpReq.req = tmp
		filterReq <- tmpReq
		// process the image process with a pipeline
		go pipelineWork(filterReq, n, done)
		<-done
	}
}

// create a pool of Readers and their work pipeline with a total number of #numOfReaders
func CreateReadersWithPPLPool(numOfReaders int, numOfThreads int) {
	var wg sync.WaitGroup
	var m sync.Mutex
	cond := sync.NewCond(&m)
	dec := json.NewDecoder(os.Stdin)

	for i := 0; i < numOfReaders; i++ {
		wg.Add(1)
		go readAndProcess(&wg, i, cond, dec, numOfThreads)
	}

	wg.Wait()
}
