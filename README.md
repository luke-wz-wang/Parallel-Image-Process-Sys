# Parallel-Image-Process-Sys


## Description


The project created a small image processing system where a series of images are read and certain effects are applied to the images using image convolution. The system reads image process requests, and the processed images are generated and output to a destinated path. The image processing system has a sequential version as well as a parallel version for handling these requests.


An image process request contains an in/source path containing the image, an out/destination path to store the generated image, and a series of effects that are to be applied to the image. A request is simulated as a single string in .json format, and Package json is used to decode the request. The system reads requests from os.Stdin and write the processed image back to the out path.


![image](https://raw.githubusercontent.com/luke-wz-wang/Post_Img/master/image-process-sys.png?token=AMHUBUEQXKJ4EVOZJH76GEC64MMXI)


These reader goroutines will read in .json format strings from os.Stdin in parallel and finish the preparation stage of the image processing. The preparation stage includes loading and initializing parts. The loading part reads the image file with .png decoder from Package image/png and instantiate an Image structure with the decoced .png image
file. The Image structure includes a in variable representing the original image, a mid pointer pointing to an image.RGBA64 instance for temporary storage in the scenario of applying multiple filter effects, and an out pointer pointing to an image.RGBA64 instance where the filtered result image is stored. The initializing part mainly instantiates a mid image from original image.


After the reader finished their preparation stage, they will call their own pipelines of workers to process the image and finally write the filtered image to the out path. Each reader has a single pipeline of workers to perform the effects. Reader goroutines will close their task channel and return when there is nothing left to read via os.Stdin from a redirect .txt file. Moreover, Package sync, specifically, sync.WaitGroup is used here to ensure the main goroutine will exit after all of the reader goroutines have returned. The communications of readers and their pipeline workers regarding the arrival of a task are achieved via channels.


The pipeline workers associated with their own readers are goroutines that actually applying the image effects on the image, and their work are handled with a pipeline pattern and data decomposition. A worker library is created for performing the work and can be found in Package src.


A pipeline is composed of a series of filtering stages (if there are multiple effects to perform) and a writer stage. Each filter stage is response for a specific filtering effect. For each filter stage, there are n_thread goroutines spawned to divide the work. The value of n_thread is read from command line via -p flag. The AddEffects() function in Package png is designed to accept an effect parameter representing a specific effect, a row index parameter y0 and a row index parameter y1, and will perform the effect on a portion of the image from row y0 to row y1. Each goroutine will perform the effect in parallel by calling the AddEffects() method with different y0 and y1 values. The first n_thread – 1 goroutines is responsible for 𝑖𝑚𝑎𝑔𝑒_h𝑒𝑖𝑔h𝑡/(𝑛_𝑡h𝑟𝑒𝑎𝑑 − 1) rows, and the last goroutine is responsible for the left part. Therefore, all goroutines will accept approximately equal portions of the image. The communications of the process status of workers of each stage and the pipeline itself are achieved via channels.


At the last stage of a pipeline, a writer goroutine will be notified via channel when the last filter stage is finished. Then, it will write the filtered image to the out path and notify the reader who owns the pipeline that the writing or entire processing stage of the image is done via channel.


## Program Usage

The program has the following optional command-line argument:

editor [-p=[number of threads]]

-p = [number of threads] = An optional flag to run the editor in its parallel version. 
