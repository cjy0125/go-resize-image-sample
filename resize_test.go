package main

import (
	"bytes"
	"crypto/md5" //nolint
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testSample = "./test/test.jpg"
)

var targetResizeFunc = ResizeImage

func testResizeImage(t *testing.T) {
	file, _ := os.Open(testSample)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		panic("cannot read from file")
	}
	file.Close()
	tmpFile := buf.Bytes()
	byteFile := make([]byte, len(tmpFile))
	copy(byteFile, tmpFile)
	thumbnail, err := targetResizeFunc(bytes.NewReader(byteFile))
	assert.Equal(t, err, nil)

	md5Sum := md5.Sum(thumbnail.Bytes()) //nolint
	t.Logf("Thumbnail md5 check sum: %x", md5Sum)
	ioutil.WriteFile(fmt.Sprintf("./debug/sample_%x.jpg", md5Sum), thumbnail.Bytes(), 0600)

	//Stress testing for multiple goroutine
	totalTestRound := 100
	results := make(chan bool, totalTestRound)
	for n := 0; n < totalTestRound; n++ {
		go func() {
			testFile := make([]byte, len(tmpFile))
			copy(testFile, tmpFile)
			smallImage, err := targetResizeFunc(bytes.NewReader(testFile))
			if err != nil {
				results <- false
				return
			}

			if md5Sum != md5.Sum(smallImage.Bytes()) { //nolint
				err := ioutil.WriteFile(fmt.Sprintf("./debug/fail_%d.jpg", time.Now().UnixNano()), smallImage.Bytes(), 0600)
				if err != nil {
					t.Log("Save file failed")
				}
				results <- false
				return
			}

			results <- true
		}()
	}

	failCount := 0
	for {
		if cap(results) == len(results) {
			t.Logf("Checking %d results", totalTestRound)
			for i := 0; i < totalTestRound; i++ {
				if !<-results {
					failCount = failCount + 1
				}
			}
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	assert.Equal(t, 0, failCount)
}

func benchmarkResizeImage(b *testing.B) {
	file, _ := os.Open(testSample)
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		panic("cannot read from file")
	}
	file.Close()
	tmpFile := buf.Bytes()
	byteFile := make([]byte, len(tmpFile))
	copy(byteFile, tmpFile)

	for n := 0; n < b.N; n++ {
		_, err := targetResizeFunc(bytes.NewReader(byteFile))
		if err != nil {
			b.Log("Encounter error in resize image")
		}
	}
}

// test
func TestBimgResizeImage(t *testing.T) {
	targetResizeFunc = bimgResizeImage
	testResizeImage(t)
}

func TestDrawResizeImage(t *testing.T) {
	targetResizeFunc = drawResizeImage
	testResizeImage(t)
}

func TestNfntResizeImage(t *testing.T) {
	targetResizeFunc = nfntResizeImage
	testResizeImage(t)
}


// benchmark
func BenchmarkBimgResizeImage(b *testing.B) {
	targetResizeFunc = bimgResizeImage
	benchmarkResizeImage(b)
}

func BenchmarkDrawResizeImage(b *testing.B) {
	targetResizeFunc = drawResizeImage
	benchmarkResizeImage(b)
}

func BenchmarkNfntResizeImage(b *testing.B) {
	targetResizeFunc = nfntResizeImage
	benchmarkResizeImage(b)
}