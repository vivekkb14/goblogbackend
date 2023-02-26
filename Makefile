.PHONY: all clean

all:
	go mod tidy
	go build -o gomysql main.go 

clean:
	-rm -rf gomysql

linux:
	go mod tidy
	GOOS=linux CGO_ENABLED=0 go build -o gomysql main.go  

windows:
	go mod tidy
	go build .

dockerise:
	go mod tidy
	GOOS=linux CGO_ENABLED=0 go build -o gomysql main.go  
	docker build -t goblog:latest .
