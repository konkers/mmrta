package mmrta

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDecodeResponse(t *testing.T) {
	file, err := os.Open("test_data/runs.json")
	if err != nil {
		t.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	_, err = decodeResponse(data)
	if err != nil {
		t.Fatal(err)
	}
}
