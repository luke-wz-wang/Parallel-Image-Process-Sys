package src

import (
	"proj2/png"
)

// create n workers, each of which is responsible for handling a part of the img
func process(img *png.Image, effect string, n int, ch chan bool){
	_, height := img.GetSize()
	// image partition
	var part int
	var last int
	if n == 1{
		part = 0
		last = height
	}else{
		// regular portion for processing
		part = height/(n-1)
		// portion left for processing
		last = height%(n-1)
	}
	// create n workers
	for i:= 0; i < n; i++{
		if i != n -1{
			go img.AddEffectRect(effect, i*part, (i+1)*part, ch)
		}else{
			// the last worker will process the left part
			go img.AddEffectRect(effect, (n-1)*part, (n-1)*part + last, ch)
		}

	}
}

// a pipeline of work stage: i.e., filter process -> ... -> filter -> process -> write out
func pipelineWork(filterReqs chan FilterRequest, n int, done chan bool) {
	request := <- filterReqs
	effects := request.req.Effects

	pipeline1 := writeStage(filterStages(request, effects, n), request)
	<-pipeline1
	done <- true
}

// write the processed image to out path
func writeStage( filterDone <- chan bool, request FilterRequest) <-chan bool {
	writeDone := make(chan bool)
	go func() {
		defer close(writeDone)
		<-filterDone
		err := request.img.Save(request.req.OutPath)
		if err != nil {
			panic(err)
		}
		writeDone <-true
	}()
	return writeDone
}

// filter stages, each of which applies a specific filtering effect
func filterStages(request FilterRequest, effects []string, n int) <-chan bool {
	filterDone := make(chan bool)
	go func() {
		defer close(filterDone)
		for i := range effects {
			ch_n := make(chan bool, n)
			process(request.img, effects[i], n, ch_n)
			for j := 0; j < n; j++ {
				<-ch_n
			}
			if i != len(effects)-1{
				request.img.Swap()
			}
		}
		filterDone <-true
	}()
	return filterDone
}
