image:
	docker build -t index-replicate:latest .

.PHONY: index-replicate
index-replicate:
	go build -o index-replicate main.go 