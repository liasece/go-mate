
.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/main github.com/liasece/go-mate

.PHONY: run
run:
	go run main/main.go buildRunner -f /home/user/testgame.go -n e1 -n e2

