package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"

	"github.com/golubaca/mjpeg"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func get(index int, url string) {

	mjpegReader := mjpeg.Reader{}
	// We can use BasicAuth as follows:
	// auth := map[string]string{"username": "cam_username", "passwd": "cam_password"}
	// err := mjpegReader.New(url, auth)
	err := mjpegReader.New(url)

	if err != nil {
		fmt.Println(err)
		runtime.Goexit()
	}

	defer mjpegReader.Close()
	fmt.Println(mjpegReader.Headers)
	for {
		img, _, _ := mjpegReader.GetFrame()

		imgName := fmt.Sprintf("%s%d%s", "/tmp/goImg", index, ".jpg")
		writeError := ioutil.WriteFile(imgName, img, 0644)
		check(writeError)

	}
}

func main() {
	// Using WaitGroup for convenience, just in case if we add more cameras
	var wg sync.WaitGroup
	cameras := map[int]string{0: "http://212.162.177.75/mjpg/video.mjpg"}
	wg.Add(len(cameras))
	fmt.Println("Total ", len(cameras), "cameras")
	for k, v := range cameras {
		go get(k, v)
		defer wg.Done()

	}

	wg.Wait()
}
