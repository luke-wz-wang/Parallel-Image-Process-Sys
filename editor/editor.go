package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"proj2/png"
	"runtime"
	"proj2/src"
)

func printUsage() {
	fmt.Printf("Usage: editor [-p=[number of threads]]\n")
	fmt.Printf("\t-p=[number of threads] = An optional flag to run the editor in its parallel version.\nCall and pass the runtime.GOMAXPROCS(...) function the integer\nspecified by [number of threads].\n")
}

func sequentialRun(){
	dec := json.NewDecoder(os.Stdin)
	var tmp src.Request

	for{
		// decode an image process request
		err := dec.Decode(&tmp)
		if err != nil {
			return
		}
		// load img
		filePath := tmp.InPath
		pngImg, err := png.Load(filePath)
		if err != nil {
			panic(err)
		}
		// apply effects
		for i := range tmp.Effects{
			_, height := pngImg.GetSize()
			pngImg.AddEffect(tmp.Effects[i], 0, height)
			if i != len(tmp.Effects) - 1{
				//pngImg.ReLoad()
				pngImg.Swap()
			}
		}
		//Saves the image to a new file
		err = pngImg.Save(tmp.OutPath)
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	/******
		The following code shows you how to work with PNG files in Golang. Please
		write your actual implementation of the project by removing the below code.
	******/
	var p int
	flag.IntVar(&p, "p", 0, "An optional flag to run the editor in its parallel version.\nCall and pass the runtime.GOMAXPROCS(...) function the integer\nspecified by [number of threads].")
	flag.Parse()

	if p == 0{
		// sequential version
		sequentialRun()
	}else{
		// parallel version
		runtime.GOMAXPROCS(p)
		numOfReaders := int(math.Ceil(float64(p)*(1.0/5.0)))
		src.CreateReadersWithPPLPool(numOfReaders, p)
	}

}
