package main

import (
	"os"
	"testing"
)

func TestCheckPath(t *testing.T) {

	path := ""

	wd, err := os.Getwd()
	if err != nil {
		t.Errorf("error: %v\n", err)
	}

	err = CheckPath(&path)
	if err != nil {
		t.Errorf("error %v\n", err)
	}

	if path != wd {
		t.Errorf("error %v\n", path)
	}
}

func TestIsTestFile(t *testing.T) {
	name := "some_test.go"
	if ok := IsTestFile(name); !ok {
		t.Errorf("error: %v\n", ok)
	}

	if ok := IsTestFile("main.go"); ok {
		t.Errorf("error: %v\n", ok)
	}
}
