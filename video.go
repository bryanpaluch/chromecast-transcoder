package main

import (
	"fmt"
)

func InspectVideo(file <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		for {
			select {
			case fileName := <-file:

				fmt.Println("Inspecting video file " + fileName)
				out <- fileName
			}
		}
	}()

	return out
}

func TranscodeVideo(file <-chan string) <-chan string {
	out := make(chan string, 4)

	go func() {
		for {
			select {
			case fileName := <-file:
				fmt.Println("Transcoding video file " + fileName)
				out <- fileName
			}
		}
	}()
	return out
}
