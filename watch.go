package main

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

func isVideoFile(name string) bool {
	ext := filepath.Ext(name)
	switch ext {
	case "mp4":
	case "mkv":
	case "avi":
		return true
	default:
		return false
	}
	return false

}

func WatchDirectory() <-chan string {
	out := make(chan string)
	go func() {
		watcher, err := NewRecursiveWatcher("./")
		if err != nil {
			log.Fatal(err)
		}
		watcher.Run(false)
		defer watcher.Close()

		for {
			select {
			case file := <-watcher.Files:

				out <- file
			}
		}
		/*

			movies := []string{"hello1.mp4", "good.mp4", "thing.mp4", "stuff.mkv"}
			out <- movies[0]
			out <- movies[1]
			/*
				for m := range movies {
					time.Sleep(1000)
					out <- movies[m]
				}
		*/

	}()

	return out
}

// From https://github.com/nathany/looper/blob/master/watch.go
//
type RecursiveWatcher struct {
	*fsnotify.Watcher
	Files   chan string
	Folders chan string
}

func NewRecursiveWatcher(path string) (*RecursiveWatcher, error) {
	folders := Subfolders(path)
	if len(folders) == 0 {
		return nil, errors.New("No folders to watch.")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	rw := &RecursiveWatcher{Watcher: watcher}

	rw.Files = make(chan string, 10)
	rw.Folders = make(chan string, len(folders))

	for _, folder := range folders {
		rw.AddFolder(folder)
	}
	return rw, nil
}

func (watcher *RecursiveWatcher) AddFolder(folder string) {
	err := watcher.Add(folder)
	if err != nil {
		log.Println("Error watching: ", folder, err)
	}
	watcher.Folders <- folder
}

func (watcher *RecursiveWatcher) Run(debug bool) {
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// create a file/directory
				if event.Op&fsnotify.Create == fsnotify.Create {
					fi, err := os.Stat(event.Name)
					if err != nil {
						// eg. stat .subl513.tmp : no such file or directory
						if debug {
							//DebugError(err)
						}
					} else if fi.IsDir() {
						if debug {
							//DebugMessage("Detected new directory %s", event.Name)
						}
						watcher.AddFolder(event.Name)
					} else {
						if debug {
							//DebugMessage("Detected new file %s", event.Name)
						}
						watcher.Files <- event.Name // created a file
					}
				}

				if event.Op&fsnotify.Write == fsnotify.Write && !(event.Op&fsnotify.Chmod == fsnotify.Chmod) {
					// modified a file, assuming that you don't modify folders
					if debug {
						//DebugMessage("Detected file modification %s", event.Name)
					}
					watcher.Files <- event.Name
				}

			case err := <-watcher.Errors:
				log.Println("error", err)
			}
		}
	}()
}

// returns a slice of subfolders (recursive), including the folder passed in
func Subfolders(path string) (paths []string) {
	filepath.Walk(path, func(newPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			// skip folders that begin with a dot
			hidden := filepath.HasPrefix(name, ".") && name != "." && name != ".."
			if hidden {
				return filepath.SkipDir
			} else {
				paths = append(paths, newPath)
			}
		}
		return nil
	})
	return paths
}
