package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting Chromecast Transcoding Service...")
	changed := WatchDirectory()
	jobs := InspectVideo(changed)
	completed := TranscodeVideo(jobs)

	<-completed
	fmt.Println("finished a job")

	<-completed
	fmt.Println("finished another job")

	<-completed
	fmt.Println("finished another job")
}
