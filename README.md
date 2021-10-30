# Go Image Resize Samples
This repository requires vips related package, therefore setup a docker image to unify the testing environment.

# Include following image resize libraries
* https://github.com/h2non/bimg
* https://github.com/nfnt/resize
* https://pkg.go.dev/golang.org/x/image/draw

# How to use
```
# build image
docker build -t go-image-resize-samples .

# invoke docker container with sh
docker run -it go-image-resize-samples /bin/sh

# run test
go test -v .

# run benchmark
go test --bench=. .
```

# Benchmark result
```
goos: linux
goarch: amd64
pkg: app
BenchmarkBimgResizeImage-2   	      27	  37308319 ns/op
BenchmarkDrawResizeImage-2   	      14	  74919800 ns/op
BenchmarkNfntResizeImage-2   	       9	 127479989 ns/op
PASS
ok  	app	22.556s
```
