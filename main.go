// testspy is a silly program that watches a directory
// and runs the "go test" command each time
// a "test" (*._test.go) file is changed.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"gopkg.in/fsnotify.v1"
)

var (
	path    = flag.String("path", "", "The directory to watch")
	watcher *fsnotify.Watcher
	done    chan bool
)

func main() {
	done := make(chan bool, 1)

	flag.Parse()

	if *path == "" {
		// Get the current working dir
		wd, err := os.Getwd()
		handleError(err)
		*path = wd
	}
	fmt.Printf("Wathing: %v\n", *path)

	matches, err := getTestFiles(*path)
	handleError(err)

	watcher, err = fsnotify.NewWatcher()
	handleError(err)
	defer watcher.Close()

	go func() {
		for {
			select {
			case ev := <-watcher.Events:

				if ev.Op&fsnotify.Create == fsnotify.Create {
					fmt.Printf("File created: %s\n", ev.Name)

					if match, _ := regexp.MatchString("\\w+_test.go", ev.Name); match {
						fmt.Printf("match: %v\n", ev.Name)

						handleError(watcher.Add(ev.Name))
						execCmd("go", "test", "./...")
					}
				}

				if ev.Op&fsnotify.Write == fsnotify.Write {
					// fmt.Printf("A file was changed: %v\n", ev)
					execCmd("go", "test", "./...")
				}
			}
		}
	}()

	addFiles(matches, *path)
	<-done
}

func getTestFiles(path string) (matches []string, err error) {

	pattern := path + "/*_test.go"

	matches, err = filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func addFiles(files []string, path string) error {
	n := len(files)
	var err error

	if n > 0 {
		for i := 0; i < n; i++ {
			if err = watcher.Add(files[i]); err != nil {
				return err
			}
		}
	}

	// if no test files are found,
	// add the working dir ...
	if path != "" {
		if err = watcher.Add(path); err != nil {
			return err
		}
	}

	return nil
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func execCmd(f string, arg ...string) {
	cmd := exec.Command(f, arg...)
	stdout, err := cmd.StdoutPipe()
	handleError(err)

	stderr, err := cmd.StderrPipe()
	handleError(err)

	err = cmd.Start()
	handleError(err)

	defer cmd.Wait()

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
}
