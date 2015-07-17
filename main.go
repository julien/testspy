package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/fsnotify.v1"
)

var (
	path       = flag.String("path", "", "The directory to watch")
	done       = make(chan bool, 1)
	fileRegexp = regexp.MustCompile("._test.go")
	watcher    *fsnotify.Watcher
)

func main() {
	flag.Parse()

	if err := checkPath(path); err != nil {
		log.Fatalf("path error: %v\n", err)
	}

	fmt.Printf("watching: %v\n", *path)

	_, err := createWatcher(*path)
	if err != nil {
		log.Fatalf("unable to create watcher: %v\n", err)
	}

	runTests()

	<-done
}

func handleError(err error) {
	if err != nil {
		fmt.Errorf("error: %v\n", err)
	}
}

func checkPath(path *string) error {
	if *path == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		*path = wd
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s\n", *path)
	}

	return nil
}

func createWatcher(path string) (watcher *fsnotify.Watcher, err error) {
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op&fsnotify.Create == fsnotify.Create || ev.Op&fsnotify.Write == fsnotify.Write {
					if isTestFile(ev.Name) {
						runTests()
					}
				}
			}
		}
	}()

	// always add base path
	if path != "" {
		if err := watcher.Add(path); err != nil {
			return nil, err
		}
	}

	return watcher, err
}

func isTestFile(name string) bool {
	return fileRegexp.MatchString(name)
}

func runTests() {
	execCmd("go", "test", "-cover", "./...")
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
