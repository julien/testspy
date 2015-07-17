package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	if *path != "" {
		log.Fatalf("error %v\n", *path)
	}
	fmt.Printf("path: %v\n", *path)

	os.Exit(m.Run())
}

func TestCheckPath(t *testing.T) {
	path := ""

	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("error: %v\n", err)
	}

	err = checkPath(&path)
	if err != nil {
		t.Errorf("error %v\n", err)
	}

	if path != wd {
		t.Errorf("error %v\n", path)
	}
}

func TestIsTestFile(t *testing.T) {
	name := "some_test.go"
	if ok := isTestFile(name); !ok {
		t.Errorf("error: %v\n", ok)
	}

	if ok := isTestFile("main.go"); ok {
		t.Errorf("error: %v\n", ok)
	}
}
