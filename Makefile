VERSION = 0.0.2

test:
	go test ./...

image:
	docker build -t index-replicate:${VERSION} .
	docker tag index-replicate:${VERSION} index-replicate:latest

.PHONY: index-replicate
index-replicate:
	go build -o index-replicate main.go

release:
	docker run --rm -it -v `pwd`:/workdir -w /workdir golang:1 make all-versions

.ONESHELL:
all-versions:
	for GOOS in darwin linux windows; do
		for GOARCH in 386 amd64; do
				echo "Building $${GOOS}-$${GOARCH}"
				export GOOS=$${GOOS}
				export GOARCH=$${GOARCH}
				go build -o bin/index-replicate-${VERSION}-$${GOOS}-$${GOARCH}
		done
	done
