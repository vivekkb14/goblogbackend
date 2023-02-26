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